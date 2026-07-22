package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

type RetrainRequest struct {
	Readings [][]float32 `json:"readings"`
}

func (apiCfg apiConfig) buildRetrainModelPayload(ctx context.Context, patientId int32) (*RetrainRequest, error) {
	duration, err := time.ParseDuration("24h")
	if err != nil {
		return nil, fmt.Errorf("Could not parse duration: %w", err)
	}
	startTime := time.Now().Add(-duration)

	readings, err := apiCfg.DB.GetReadings(
		ctx,
		database.GetReadingsParams{
			PatientID:     patientId,
			TimeOfReading: startTime,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Could not get readings: %w", err)
	}

	modelInput := make([][]float32, 0)
	for _, reading := range readings {
		gsm, err := strconv.ParseFloat(reading.Glucose, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid glucose value: %w", err)
		}

		basalRate, err := strconv.ParseFloat(reading.BasalRate, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid basal_rate value: %w", err)
		}

		bolus, err := strconv.ParseFloat(reading.Bolus, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid bolus value: %w", err)
		}

		insulin := basalRate + bolus

		meal, err := strconv.ParseFloat(reading.Carbs, 32)
		if err != nil {
			return nil, fmt.Errorf("Invalid carbs value: %w", err)
		}

		modelInput = append(modelInput, []float32{
			float32(gsm),
			float32(insulin),
			float32(meal),
			float32(reading.ExerciseDuration.Int32),
			float32(reading.ExerciseIntensity.Int32),
		})
	}

	return &RetrainRequest{
		Readings: modelInput,
	}, nil
}

func (apiCfg apiConfig) sendRetrainPayload(ctx context.Context, patientID int32, modelID uuid.UUID, signal string, payload *RetrainRequest) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Could not marshal retrain payload: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		apiCfg.ModelURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("Could not build retrain request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("patient_id", strconv.FormatInt(int64(patientID), 10))
	req.Header.Set("model_id", modelID.String())
	req.Header.Set("signal", signal)

	resp, err := apiCfg.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Could not send retrain request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Retrain service returned status %d", resp.StatusCode)
	}

	return nil
}

func (apiCfg apiConfig) handlerSendInforForRetraining(w http.ResponseWriter, r *http.Request, patientID int32, modelID uuid.UUID) error {
	payload, err := apiCfg.buildRetrainModelPayload(r.Context(), patientID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not build retraining payload: %s", err))
		return err
	}

	err = apiCfg.sendRetrainPayload(r.Context(), patientID, modelID, "retrain", payload)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not send recommendations payload: %s", err))
		return err
	}
	err = apiCfg.sendRetrainPayload(r.Context(), patientID, modelID, "recalibrate", payload)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not send recommendations payload: %s", err))
		return err
	}

	return nil
}
