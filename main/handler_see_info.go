package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
)

func (apiCfg *apiConfig) handlerSeeInfo(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Readings         []database.Reading                `json:"readings"`
		ModelPredictions []database.GetModelPredictionsRow `json:"model_predictions"`
		OduPredictions   []database.GetModelPredictionsRow `json:"odu_predictions"`
	}

	type parameters struct {
		PatientID int32 `json:"patient"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return

	}

	readings, err := apiCfg.DB.GetReadings(
		r.Context(),
		params.PatientID,
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get readings: %s", err))
		return
	}

	now := time.Now()

	modelPredictions, err := apiCfg.DB.GetModelPredictions(
		r.Context(),
		database.GetModelPredictionsParams{
			PatientID:     params.PatientID,
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
			PatientID:     params.PatientID,
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
