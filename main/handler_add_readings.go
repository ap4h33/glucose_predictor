package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerAddReadings(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Patient           uuid.NullUUID `json:"patient"`
		TimeOfReading     time.Time     `json:"timestamp"`
		Glucose           string        `json:"glucose"`
		BasalRate         string        `json:"basal_rate"`
		Bolus             string        `json:"bolus"`
		Carbs             string        `json:"carbs"`
		ExerciseDuration  int32         `json:"ex_duration"`
		ExerciseIntensity int32         `json:"ex_intensity"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	ExDuration := sql.NullInt32{
		Int32: params.ExerciseDuration,
		Valid: true,
	}
	ExIntensity := sql.NullInt32{
		Int32: int32(params.ExerciseIntensity),
		Valid: true,
	}

	reading, err := apiCfg.DB.AddReading(r.Context(), database.AddReadingParams{
		ID:                uuid.New(),
		PatientID:         params.Patient,
		TimeOfReading:     params.TimeOfReading,
		Glucose:           params.Glucose,
		BasalRate:         params.BasalRate,
		Bolus:             params.Bolus,
		Carbs:             params.Carbs,
		ExerciseDuration:  ExDuration,
		ExerciseIntensity: ExIntensity,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create a reading: %s", err))
	}
	respondWithJSON(w, 200, reading)

}
