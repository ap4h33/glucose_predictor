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
	//This function sends over patient history and the latest predictons for glucose levels
	type response struct {
		Readings         []database.Reading                `json:"readings"`
		ModelPredictions []database.GetModelPredictionsRow `json:"model_predictions"`
		OduPredictions   []database.GetModelPredictionsRow `json:"odu_predictions"`
	}

	patientID, err := strconv.ParseInt(chi.URLParam(r, "patient_id"), 10, 32)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %s", err))
		return
	}

	period := r.URL.Query().Get("time_period")
	var startTime time.Time
	getAll := false

	if period == "" {
		period = "24h"
	}
	if period == "all" {
		getAll = true
	} else {
		duration, err := time.ParseDuration(period)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Invalid time period: %s", err))
			return
		}
		startTime = time.Now().Add(-duration)
	}

	var readings []database.Reading

	if getAll {
		readings, err = apiCfg.DB.GetAllReadings(
			r.Context(),
			int32(patientID),
		)
	} else {
		readings, err = apiCfg.DB.GetReadings(
			r.Context(),
			database.GetReadingsParams{
				PatientID:     int32(patientID),
				TimeOfReading: startTime,
			},
		)
	}
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get readings: %s", err))
		return
	}

	now := time.Now()

	aiPredictions, err := apiCfg.DB.GetModelPredictions(
		r.Context(),
		database.GetModelPredictionsParams{
			Name:          "ai_model",
			Version:       apiCfg.AIVersion,
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
			Name:          "odu_model",
			Version:       apiCfg.ODUVersion,
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
		ModelPredictions: aiPredictions,
		OduPredictions:   oduPredictions,
	})
}
