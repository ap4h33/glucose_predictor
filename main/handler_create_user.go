package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ap4h33/glucose_predictor/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	reading, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Column1: int32(params.ID),
		Name:    params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user: %s", err))
	}
	respondWithJSON(w, 200, reading)

}
