package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
)

func (apiCfg apiConfig) handlerGetRecommendations(w http.ResponseWriter, r *http.Request) {

	patient_id, err := strconv.ParseInt(chi.URLParam(r, "patient_id"), 10, 32)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %s", err))
		return
	}

	recommendation, err := apiCfg.DB.GetRecommendation(
		r.Context(),
		int32(patient_id),
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get a recommendation: %s", err))
		return
	}

	respondWithJSON(w, 200, recommendation)
}

func (apiCfg apiConfig) handlerSendInforForRecommendations(w http.ResponseWriter, r *http.Request, patientID int32) {
	type response struct {
		Readings         []database.Reading                `json:"readings"`
		ModelPredictions []database.GetModelPredictionsRow `json:"model_predictions"`
		OduPredictions   []database.GetModelPredictionsRow `json:"odu_predictions"`
	}

	type PredictionResponse struct {
		ModelPrediction database.GetModelPredictionsRow
		ModelName       string `json:"model_name"`
	}

	now := time.Now()

	// AIModel :=
	// oduModel :=

	aiPredictions, err := apiCfg.DB.GetModelPredictions(
		r.Context(),
		database.GetModelPredictionsParams{
			// ModelID: AIModel,
			PatientID:     patientID,
			TimePredicted: now,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get model predictions: %s", err))
		return
	}

	oduPredictions, err := apiCfg.DB.GetModelPredictions(
		r.Context(),
		database.GetModelPredictionsParams{
			// ModelID: oduModel,
			PatientID:     patientID,
			TimePredicted: now,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get model predictions: %s", err))
		return
	}

	var predictions []PredictionResponse
	// This creates an array of all of the predictions, and labals them with according names.
	for _, p := range aiPredictions {
		predictions = append(predictions, PredictionResponse{
			ModelPrediction: p,
			ModelName:       "AI",
		})
	}

	for _, p := range oduPredictions {
		predictions = append(predictions, PredictionResponse{
			ModelPrediction: p,
			ModelName:       "ODU",
		})
	}

}

// func (apiCfg apiConfig) handlerCreateRecommendations(w http.ResponseWriter, r *http.Request) {

// }
