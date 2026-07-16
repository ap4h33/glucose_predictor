package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID         int           `json:"id"`
		Name       string        `json:"name"`
		HospitalID uuid.NullUUID `json:"hospital_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	reading, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Column1:    int32(params.ID),
		HospitalID: params.HospitalID,
		Name:       params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %s", err))
	}
	respondWithJSON(w, 200, reading)

}

func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(chi.URLParam(r, "patient_id"), 10, 32)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %s", err))
		return
	}

	user, err := apiCfg.DB.GetUser(
		r.Context(),
		int32(id),
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get user: %s", err))
		return
	}

	respondWithJSON(w, 200, user)
}

func (apiCfg apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {

	adminID, err := uuid.Parse(chi.URLParam(r, "admin_id"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid admin ID: %s", err))
		return
	}

	users, err := apiCfg.DB.GetUsers(
		r.Context(),
		adminID,
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get users: %s", err))
		return
	}

	respondWithJSON(w, 200, users)
}
