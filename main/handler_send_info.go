package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerSendUnseenReadingsToModel(patientID uuid.UUID) error {
	//change this value to control how many readings are needed before sending
	const requiredReadings = 5

	readings, err := apiCfg.DB.GetUnseenReadings(
		context.Background(),
		patientID,
	)
	if err != nil {
		return fmt.Errorf("could not get unseen readings: %s", err)
	}

	if len(readings) < requiredReadings {
		return nil
	}

	type modelReading struct {
		Patient           uuid.UUID `json:"patient"`
		TimeOfReading     time.Time `json:"timestamp"`
		Glucose           string    `json:"glucose"`
		BasalRate         string    `json:"basal_rate"`
		Bolus             string    `json:"bolus"`
		Carbs             string    `json:"carbs"`
		ExerciseDuration  int32     `json:"ex_duration"`
		ExerciseIntensity int32     `json:"ex_intensity"`
	}

	modelReadings := make([]modelReading, 0, len(readings))

	for _, reading := range readings {
		modelReadings = append(modelReadings, modelReading{
			Patient:           reading.PatientID,
			TimeOfReading:     reading.TimeOfReading,
			Glucose:           reading.Glucose,
			BasalRate:         reading.BasalRate,
			Bolus:             reading.Bolus,
			Carbs:             reading.Carbs,
			ExerciseDuration:  reading.ExerciseDuration.Int32,
			ExerciseIntensity: reading.ExerciseIntensity.Int32,
		})
	}

	body, err := json.Marshal(modelReadings)
	if err != nil {
		return fmt.Errorf("could not encode readings: %s", err)
	}

	response, err := http.Post(
		apiCfg.ModelURL,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("could not send readings to model: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("model returned status code: %d", response.StatusCode)
	}

	for _, reading := range readings {
		err := apiCfg.DB.SendReadingToModel(
			context.Background(),
			reading.ID,
		)
		if err != nil {
			return fmt.Errorf("could not update reading %s: %s", reading.ID, err)
		}
	}

	return nil
}
