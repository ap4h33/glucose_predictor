package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerAddPredictions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ModelID          uuid.UUID   `json:"model_id"`
		PatientID        int32       `json:"patient_id"`
		GlucosePredicted json.Number `json:"glucose_predicted"`
		TimePredicted    time.Time   `json:"time_predicted"`
		GeneratedAt      time.Time   `json:"generated_at"`
	}

	var params []parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	predictions := make([]database.Prediction, 0, len(params))

	for _, p := range params {
		prediction, err := apiCfg.DB.AddPrediction(r.Context(), database.AddPredictionParams{
			ID:               uuid.New(),
			ModelID:          p.ModelID,
			PatientID:        p.PatientID,
			GlucosePredicted: p.GlucosePredicted.String(),
			TimePredicted:    p.TimePredicted,
			GeneratedAt:      p.GeneratedAt,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Could not create prediction: %s", err))
			return
		}

		predictions = append(predictions, prediction)
	}

	respondWithJSON(w, http.StatusOK, predictions)
}
