package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ap4h33/glucose_predictor/internal/database"
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
	type parameters struct {
		UserID int32 `json:"patient_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse JSON: %s", err))
		return
	}

	user, err := apiCfg.DB.GetUser(
		r.Context(),
		params.UserID,
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get user: %s", err))
		return
	}

	respondWithJSON(w, 200, user)
}

func (apiCfg apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		AdminID uuid.UUID `json:"admin_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not parse JSON: %s", err))
	}

	users, err := apiCfg.DB.GetUsers(
		r.Context(),
		params.AdminID,
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get users: %s", err))
		return
	}

	respondWithJSON(w, 200, users)
}
