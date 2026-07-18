package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (apiCfg *apiConfig) handlerSendUnseenReadingsToModel(patientID int32) error {
	const readings_amount = 12
	readings, err := apiCfg.DB.GetLastReadings(
		context.Background(),
		patientID,
	)
	if err != nil {
		return fmt.Errorf("could not get unseen readings: %w", err)
	}

	// Not enough data for the model
	if len(readings) < readings_amount {
		return nil
	}

	// Shape: (1, readings_amount, 6)
	modelInput := make([][][]float32, 1)
	modelInput[0] = make([][]float32, 0, len(readings))

	for _, reading := range readings {
		glucose, err := strconv.ParseFloat(reading.Glucose, 32)
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

		carbs, err := strconv.ParseFloat(reading.Carbs, 32)
		if err != nil {
			return fmt.Errorf("invalid carbs value: %w", err)
		}

		modelInput[0] = append(modelInput[0], []float32{
			float32(glucose),
			float32(basalRate),
			float32(bolus),
			float32(carbs),
			float32(reading.ExerciseDuration.Int32),
			float32(reading.ExerciseIntensity.Int32),
		})
	}

	body, err := json.Marshal(modelInput)
	if err != nil {
		return fmt.Errorf("could not encode readings: %w", err)
	}

	response, err := http.Post(
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

	// Marks the readings that were sent
	for _, reading := range readings {
		err := apiCfg.DB.UpdateReadingModelStatus(
			context.Background(),
			reading.ID,
		)
		if err != nil {
			return fmt.Errorf("could not update reading %s: %w", reading.ID, err)
		}
	}

	return nil
}
