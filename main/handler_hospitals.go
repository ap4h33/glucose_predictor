package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateHospital(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		HospitalName string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not decode JSON: %s", err))
		return
	}

	hospital, err := apiCfg.DB.CreateHospital(r.Context(), database.CreateHospitalParams{
		ID:   uuid.New(),
		Name: params.HospitalName,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create hospital: %s", err))
		return
	}
	respondWithJSON(w, 200, hospital)
}

func (apiCfg apiConfig) handlerGetHospitals(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "hospital_id")

	if idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondWithError(w, 400, "Invalid hospital ID")
			return
		}
		hospital, err := apiCfg.DB.GetHospital(
			r.Context(),
			id,
		)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Could not get hospital: %s", err))
			return
		}
		respondWithJSON(w, 200, hospital)
		return
	}

	hospitals, err := apiCfg.DB.GetAllHospitals(
		r.Context(),
	)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get hospitals: %s", err))
		return
	}

	respondWithJSON(w, 200, hospitals)
}
