package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateAdmin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		HospitalID uuid.NullUUID `json:"hospital_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	reading, err := apiCfg.DB.CreateAdmin(r.Context(), database.CreateAdminParams{
		ID:         uuid.New(),
		HospitalID: params.HospitalID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %s", err))
	}
	respondWithJSON(w, 200, reading)

}

func (apiCfg apiConfig) handlerGetAdmin(w http.ResponseWriter, r *http.Request) {

	adminID, err := uuid.Parse(chi.URLParam(r, "admin_id"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid admin ID: %s", err))
		return
	}

	admin, err := apiCfg.DB.GetAdmin(
		r.Context(),
		adminID,
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get admin: %s", err))
		return
	}

	respondWithJSON(w, 200, admin)
}
