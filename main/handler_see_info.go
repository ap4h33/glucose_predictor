package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
)

func (apiCfg *apiConfig) handlerSeeInfo(w http.ResponseWriter, r *http.Request) {
	//This function send over full patient history and the latest predictons for glucose levels
	type response struct {
		Readings         []database.Reading                `json:"readings"`
		ModelPredictions []database.GetModelPredictionsRow `json:"model_predictions"`
		OduPredictions   []database.GetModelPredictionsRow `json:"odu_predictions"`
	}

	patientID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 32)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %s", err))
		return
	}

	readings, err := apiCfg.DB.GetReadings(
		r.Context(),
		int32(patientID),
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get readings: %s", err))
		return
	}

	now := time.Now()

	modelPredictions, err := apiCfg.DB.GetModelPredictions(
		r.Context(),
		database.GetModelPredictionsParams{
			PatientID:     int32(patientID),
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
			PatientID:     int32(patientID),
			TimePredicted: now,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get ODU predictions: %s", err))
		return
	}

	respondWithJSON(w, 200, response{
		Readings:         readings,
		ModelPredictions: modelPredictions,
		OduPredictions:   oduPredictions,
	})
}
