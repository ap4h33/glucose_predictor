package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerAddReadings(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Patient           int32       `json:"patient_id"`
		TimeOfReading     time.Time   `json:"timestamp"`
		Glucose           json.Number `json:"glucose"`
		BasalRate         json.Number `json:"basal_rate"`
		Bolus             json.Number `json:"bolus"`
		Carbs             json.Number `json:"carbs"`
		ExerciseDuration  int32       `json:"exercise_duration"`
		ExerciseIntensity int32       `json:"exercise_intensity"`
	}

	var params []parameters

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	readings := make([]database.Reading, 0, len(params))
	patients := make(map[int32]struct{})
	//makes an empty slice, appends processed reading into said slice, then returns the whole thing once the loop is complete
	for _, p := range params {
		reading, err := apiCfg.DB.AddReading(r.Context(), database.AddReadingParams{
			ID:            uuid.New(),
			PatientID:     p.Patient,
			TimeOfReading: p.TimeOfReading,
			Glucose:       p.Glucose.String(),
			BasalRate:     p.BasalRate.String(),
			Bolus:         p.Bolus.String(),
			Carbs:         p.Carbs.String(),
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
		patients[reading.PatientID] = struct{}{}
	}

	for patientID := range patients {
		go func(id int32) {
			if err := apiCfg.handlerSendReadingsToModel(id); err != nil {
				log.Printf("error sending readings for patient %d: %v", id, err)
			}
		}(patientID)
	}

	respondWithJSON(w, 200, readings)
}
