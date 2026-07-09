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

	var params []parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	readings := make([]database.Reading, 0, len(params))
	//makes an epty slice, appends processed reading into said slice, then returns the whole thing once the loop is complete
	for _, p := range params {
		reading, err := apiCfg.DB.AddReading(r.Context(), database.AddReadingParams{
			ID:            uuid.New(),
			PatientID:     p.Patient,
			TimeOfReading: p.TimeOfReading,
			Glucose:       p.Glucose,
			BasalRate:     p.BasalRate,
			Bolus:         p.Bolus,
			Carbs:         p.Carbs,
			ExerciseDuration: sql.NullInt32{
				Int32: p.ExerciseDuration,
				Valid: true,
			},
			ExerciseIntensity: sql.NullInt32{
				Int32: p.ExerciseIntensity,
				Valid: true,
			},
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Could not create reading: %s", err))
			return
		}

		readings = append(readings, reading)
	}

	respondWithJSON(w, 200, readings)
}
