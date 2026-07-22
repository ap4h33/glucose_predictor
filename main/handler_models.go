package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ap4h33/glucose_predictor/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerAddModel(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	model, err := apiCfg.DB.AddModel(r.Context(), database.AddModelParams{
		ID:      uuid.New(),
		Name:    params.Name,
		Version: params.Version,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create model: %s", err))
		return
	}

	respondWithJSON(w, 200, model)
}

func (apiCfg *apiConfig) handlerGetAllModels(w http.ResponseWriter, r *http.Request) {
	models, err := apiCfg.DB.GetAllModels(
		r.Context(),
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get models: %s", err))
		return
	}

	respondWithJSON(w, 200, models)
}
