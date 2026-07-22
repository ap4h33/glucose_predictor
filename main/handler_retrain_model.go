package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
)

type RetrainRequest struct {
	Signal   string      // TO DO: MOVE TO HEADER
	Readings [][]float32 `json:"readings"`
}

func (apiCfg apiConfig) handlerBuildRetrainModelPayload(patientId int32) (*RetrainRequest, error) {

	// TO DO: ADD PATIENT ID AND MODEL ID IN HEADER

	duration, err := time.ParseDuration("24h")
	if err != nil {
		return nil, fmt.Errorf("Could not parse duration: %w", err)
	}
	startTime := time.Now().Add(-duration)

	readings, err := apiCfg.DB.GetReadings(
		r.Context(),
		database.GetReadingsParams{
			PatientID:     patientId,
			TimeOfReading: startTime,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Cuold not get readings: %w", err)
	}

	modelInput := make([][]float32, 0)
	for _, reading := range readings {
		gsm, err := strconv.ParseFloat(reading.Glucose, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid glucose value: %w", err)
		}

		basalRate, err := strconv.ParseFloat(reading.BasalRate, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid basal_rate value: %w", err)
		}

		bolus, err := strconv.ParseFloat(reading.Bolus, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid bolus value: %w", err)
		}

		insulin := basalRate + bolus

		meal, err := strconv.ParseFloat(reading.Carbs, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid carbs value: %w", err)
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
