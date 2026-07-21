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
		prediction30, err := apiCfg.DB.AddPrediction(r.Context(), database.AddPredictionParams{
			ID:               uuid.New(),
			ModelID:          modelID,
			PatientID:        patientID,
			GlucosePredicted: p.GlucoseIn30.String(),
			TimePredicted:    lastReadingTime.Add(30 * time.Minute),
			GeneratedAt:      createdAt,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create prediction: %s", err))
			return
		}

		prediction60, err := apiCfg.DB.AddPrediction(r.Context(), database.AddPredictionParams{
			ID:               uuid.New(),
			ModelID:          modelID,
			PatientID:        patientID,
			GlucosePredicted: p.GlucoseIn60.String(),
			TimePredicted:    lastReadingTime.Add(60 * time.Minute),
			GeneratedAt:      createdAt,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create prediction: %s", err))
			return
		}

		prediction90, err := apiCfg.DB.AddPrediction(r.Context(), database.AddPredictionParams{
			ID:               uuid.New(),
			ModelID:          modelID,
			PatientID:        patientID,
			GlucosePredicted: p.GlucoseIn90.String(),
			TimePredicted:    lastReadingTime.Add(90 * time.Minute),
			GeneratedAt:      createdAt,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create prediction: %s", err))
			return
		}

		predictions = append(predictions, prediction30)
		predictions = append(predictions, prediction60)
		predictions = append(predictions, prediction90)
	}

	err = apiCfg.handlerSendInforForRecommendations(w, r, patientID)
	if err != nil {
		return
	}

	respondWithJSON(w, 200, predictions)
}
