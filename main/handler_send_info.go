package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	response, err := http.Post(
		// TO DO: differenciaet models
		apiCfg.ModelURL,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("could not send readings to model: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("model returned status code: %d", response.StatusCode)
	}

	return nil
}
