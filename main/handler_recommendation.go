package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

type ReadingsResponse struct {
	Glucose           string `json:"glucose"`
	BasalRate         string `json:"basal_rate"`
	Bolus             string `json:"bolus"`
	Carbs             string `json:"carbs"`
	ExerciseDuration  int32  `json:"exercise_duration"`
	ExerciseIntensity int32  `json:"exercise_intensity"`
}

type PredictionsResponse struct {
	ModelPrediction database.GetModelPredictionsRow
	ModelName       string `json:"model_name"`
}

type RecommendationsRequest struct {
	Readings    []ReadingsResponse    `json:"readings"`
	Predictions []PredictionsResponse `json:"predictions"`
}

func (apiCfg apiConfig) buildRecommendationsPayload(ctx context.Context, patientID int32) (RecommendationsRequest, time.Time, error) {
	// Time is used to get predictions > now, and readings from 1h ago till now.
	now := time.Now()
	duration, err := time.ParseDuration("1h")
	if err != nil {
		return RecommendationsRequest{}, time.Time{}, fmt.Errorf("could not parse time: %w", err)
	}
	readingsTime := now.Add(-duration)

	readings, err := apiCfg.DB.GetReadings(
		ctx,
		database.GetReadingsParams{
			PatientID:     patientID,
			TimeOfReading: readingsTime,
		},
	)
	if err != nil {
		return RecommendationsRequest{}, time.Time{}, fmt.Errorf("could not get readings: %w", err)
	}

	var lastReadingTime time.Time
	if len(readings) > 0 {
		lastReadingTime = readings[0].TimeOfReading
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

	aiPredictions, err := apiCfg.DB.GetModelPredictions(
		ctx,
		database.GetModelPredictionsParams{
			Name:          "ai_model",
			Version:       apiCfg.AIVersion,
			PatientID:     patientID,
			TimePredicted: now,
		},
	)
	if err != nil {
		return RecommendationsRequest{}, time.Time{}, fmt.Errorf("could not get model predictions: %w", err)
	}

	oduPredictions, err := apiCfg.DB.GetModelPredictions(
		ctx,
		database.GetModelPredictionsParams{
			Name:          "odu_model",
			Version:       apiCfg.ODUVersion,
			PatientID:     patientID,
			TimePredicted: now,
		},
	)
	if err != nil {
		return RecommendationsRequest{}, time.Time{}, fmt.Errorf("could not get model predictions: %w", err)
	}

	var predictionsResponse []PredictionsResponse
	for _, p := range aiPredictions {
		predictionsResponse = append(predictionsResponse, PredictionsResponse{
			ModelPrediction: p,
			ModelName:       "AI",
		})
	}
	for _, p := range oduPredictions {
		predictionsResponse = append(predictionsResponse, PredictionsResponse{
			ModelPrediction: p,
			ModelName:       "ODU",
		})
	}

	return RecommendationsRequest{
		Readings:    readingsResponse,
		Predictions: predictionsResponse,
	}, lastReadingTime, nil
}

func (apiCfg apiConfig) sendRecommendationsPayload(ctx context.Context, patientID int32, lastReadingTime time.Time, payload RecommendationsRequest) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal recommendations payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiCfg.RecURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not build recommendations request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("patient_id", strconv.FormatInt(int64(patientID), 10))
	if !lastReadingTime.IsZero() {
		req.Header.Set("last_reading_time", lastReadingTime.Format(time.RFC3339))
	}

	resp, err := apiCfg.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send info for recommendations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("recommendations service returned status %d", resp.StatusCode)
	}

	return nil
}

func (apiCfg apiConfig) handlerSendInforForRecommendations(w http.ResponseWriter, r *http.Request, patientID int32) error {
	payload, lastReadingTime, err := apiCfg.buildRecommendationsPayload(r.Context(), patientID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not build recommendations payload: %s", err))
		return err
	}

	err = apiCfg.sendRecommendationsPayload(r.Context(), patientID, lastReadingTime, payload)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not send recommendations payload: %s", err))
		return err
	}

	return nil
}

func (apiCfg apiConfig) handlerAddRecommendations(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Patient   int32     `json:"patient_id"`
		IsSafe    bool      `json:"is_safe"`
		Action    string    `json:"action"`
		Message   string    `json:"message"`
		CreatedAt time.Time `json:"created_at"`
	}

	var params []parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	recommendations := make([]database.Recommendation, 0, len(params))
	//makes an empty slice, appends processed recommendation into said slice, then returns the whole thing once the loop is complete
	for _, p := range params {
		recommendation, err := apiCfg.DB.AddRecommendation(r.Context(), database.AddRecommendationParams{
			ID:        uuid.New(),
			PatientID: p.Patient,
			IsSafe:    p.IsSafe,
			Action: sql.NullString{
				String: p.Action,
				Valid:  p.Action != "",
			},
			Message: sql.NullString{
				String: p.Message,
				Valid:  p.Message != "",
			},
			CreatedAt: sql.NullTime{
				Time:  p.CreatedAt,
				Valid: !p.CreatedAt.IsZero(),
			},
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Could not create recommendation: %s", err))
			return
		}

		recommendations = append(recommendations, recommendation)
	}

	respondWithJSON(w, 200, recommendations)
}
