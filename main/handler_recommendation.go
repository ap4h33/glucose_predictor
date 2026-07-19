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
	// Time is used to get predictions > now, and readings from 1h ago till now.
	now := time.Now()
	duration, err := time.ParseDuration("1h")
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse time: %s", err))
	}
	readingsTime := now.Add(-duration)

	// Readings
	type ReadingsResponse struct {
		Glucose           string `json:"glucose"`
		BasalRate         string `json:"basal_rate"`
		Bolus             string `json:"bolus"`
		Carbs             string `json:"carbs"`
		ExerciseDuration  int32  `json:"exercise_duration"`
		ExerciseIntensity int32  `json:"exercise_intensity"`
	}

	readings, err := apiCfg.DB.GetReadings(
		r.Context(),
		database.GetReadingsParams{
			PatientID:     patientID,
			TimeOfReading: readingsTime,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get readings: %s", err))
	}

	if len(readings) > 0 {
		w.Header().Set(
			"parient_id",
			strconv.FormatInt(int64(patientID), 10),
		)
		w.Header().Set(
			"last_reading_time",
			readings[0].TimeOfReading.Format(time.RFC3339),
		)
	}

	for i, j := 0, len(readings)-1; i < j; i, j = i+1, j-1 {
		readings[i], readings[j] = readings[j], readings[i]
	}

	var readingsResponse []ReadingsResponse
	for _, read := range readings {
		readingsResponse = append(readingsResponse, ReadingsResponse{
			Glucose:           read.Glucose,
			BasalRate:         read.BasalRate,
			Bolus:             read.Bolus,
			Carbs:             read.Carbs,
			ExerciseDuration:  read.ExerciseDuration.Int32,
			ExerciseIntensity: read.ExerciseIntensity.Int32,
		})
	}

	// Predictions
	type PredictionResponse struct {
		ModelPrediction database.GetModelPredictionsRow
		ModelName       string `json:"model_name"`
	}

	aiPredictions, err := apiCfg.DB.GetModelPredictions(
		r.Context(),
		database.GetModelPredictionsParams{
			Name:          "ai_model",
			Version:       apiCfg.AIVersion,
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
			Name:          "odu_model",
			Version:       apiCfg.ODUVerion,
			PatientID:     int32(patientID),
			TimePredicted: now,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get model predictions: %s", err))
		return
	}

	var predictionsResponse []PredictionResponse
	// This creates an array of all of the predictions, and labals them with according names.
	for _, p := range aiPredictions {
		predictionsResponse = append(predictionsResponse, PredictionResponse{
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

	type response struct {
		Readings    []ReadingsResponse   `json:"readings"`
		Predictions []PredictionResponse `json:"predictions"`
	}

	respondWithJSON(w, 200, response{
		Readings:    readingsResponse,
		Predictions: predictionsResponse,
	})

}

// func (apiCfg apiConfig) handlerCreateRecommendations(w http.ResponseWriter, r *http.Request) {

// }
