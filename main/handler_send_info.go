package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerSendReadingsToModel(patientID int32) error {
	const readings_amount = 12
	type ModelReading struct {
		Timestamp         string  `json:"timestamp"`
		Glucose           float32 `json:"glucose"`
		Insulin           float32 `json:"insulin"`
		Meal              float32 `json:"meal"`
		ExerciseDuration  float32 `json:"exercise_duration"`
		ExerciseIntensity float32 `json:"exercise_intensity"`
	}

	type ModelRequest struct {
		Readings []ModelReading `json:"readings"`
	}
	readings, err := apiCfg.DB.GetLastReadings(
		context.Background(),
		patientID,
	)
	if err != nil {
		return fmt.Errorf("could not get readings: %w", err)
	}

	// Для новых пользователей, если недостаточно записей с девайса, предсказания не запускаются
	if len(readings) < 12 {
		return nil
	}
	// Shape: (1, readings_amount, 6)
	modelInput := make([][][]float32, 1)
	modelInput[0] = make([][]float32, 0, len(readings))

	for _, reading := range readings {
		gsm, err := strconv.ParseFloat(reading.Glucose, 32)
		if err != nil {
			return fmt.Errorf("invalid glucose value: %w", err)
		}

		basalRate, err := strconv.ParseFloat(reading.BasalRate, 32)
		if err != nil {
			return fmt.Errorf("invalid basal_rate value: %w", err)
		}

		bolus, err := strconv.ParseFloat(reading.Bolus, 32)
		if err != nil {
			return fmt.Errorf("invalid bolus value: %w", err)
		}

		insulin := basalRate + bolus

		meal, err := strconv.ParseFloat(reading.Carbs, 32)
		if err != nil {
			return fmt.Errorf("invalid carbs value: %w", err)
		}

		modelInput[0] = append(modelInput[0], []float32{
			float32(gsm),
			float32(insulin),
			float32(meal),
			float32(reading.ExerciseDuration.Int32),
			float32(reading.ExerciseIntensity.Int32),
		})
	}

	body, err := json.Marshal(modelInput)
	if err != nil {
		return fmt.Errorf("could not encode readings: %w", err)
	}

	client := &http.Client{}

	modelIDs := []uuid.UUID{
		apiCfg.AImodelID,
		apiCfg.ODUmodelID,
	}

	for _, modelID := range modelIDs {
		req, err := http.NewRequest(
			http.MethodPost,
			apiCfg.ModelURL,
			bytes.NewReader(body),
		)
		if err != nil {
			return fmt.Errorf("could not create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Patient-ID", strconv.Itoa(int(patientID)))
		req.Header.Set("Model-ID", modelID.String())

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("could not send readings to model %s: %w", modelID, err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return fmt.Errorf("model %s returned status code: %d", modelID, resp.StatusCode)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return fmt.Errorf("model %s returned status code: %d", modelID, resp.StatusCode)
		}

		resp.Body.Close()
	}

	return nil
}
