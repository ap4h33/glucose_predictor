package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerAddPredictions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GlucoseIn30 json.Number `json:"glucose30"`
		GlucoseIn60 json.Number `json:"glucose60"`
		GlucoseIn90 json.Number `json:"glucose90"`
	}

	var params []parameters

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	modelID, err := uuid.Parse(r.Header.Get("Model-ID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Model-ID header")
		return
	}

	patientID64, err := strconv.ParseInt(r.Header.Get("Patient-ID"), 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Patient-ID header")
		return
	}
	patientID := int32(patientID64)

	lastReadingTime, err := time.Parse(time.RFC3339, r.Header.Get("Last-Reading-Time"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Last-Reading-Time header")
		return
	}

	createdAt, err := time.Parse(time.RFC3339, r.Header.Get("Created-At"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Created-At header")
		return
	}

	predictions := make([]database.Prediction, 0, len(params)*3)

	for _, p := range params {
		predictionData := []struct {
			glucose json.Number
			minutes time.Duration
		}{
			{p.GlucoseIn30, 30 * time.Minute},
			{p.GlucoseIn60, 60 * time.Minute},
			{p.GlucoseIn90, 90 * time.Minute},
		}

		for _, pred := range predictionData {
			prediction, err := apiCfg.DB.AddPrediction(r.Context(), database.AddPredictionParams{
				ID:               uuid.New(),
				ModelID:          modelID,
				PatientID:        patientID,
				GlucosePredicted: pred.glucose.String(),
				TimePredicted:    lastReadingTime.Add(pred.minutes),
				GeneratedAt:      createdAt,
			})
			if err != nil {
				respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create prediction: %s", err))
				return
			}

			predictions = append(predictions, prediction)
		}
	}

	err = apiCfg.handlerSendInforForRecommendations(w, r, patientID)
	if err != nil {
		return
	}

	respondWithJSON(w, 200, predictions)
}

func (apiCfg *apiConfig) handlerGetAllModelPredictions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		PatientID int32  `json:"patient_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	predictions, err := apiCfg.DB.GetAllModelPredictions(
		r.Context(),
		database.GetAllModelPredictionsParams{
			Name:      params.Name,
			Version:   params.Version,
			PatientID: params.PatientID,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get model predictions: %s", err))
		return
	}

	respondWithJSON(w, 200, predictions)
}
