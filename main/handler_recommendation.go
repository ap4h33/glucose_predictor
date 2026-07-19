package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (apiCfg apiConfig) handlerGetRecommendations(w http.ResponseWriter, r *http.Request) {

	patient_id, err := strconv.ParseInt(chi.URLParam(r, "patient_id"), 10, 32)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid ID: %s", err))
		return
	}

	recommendation, err := apiCfg.DB.GetRecommendation(
		r.Context(),
		int32(patient_id),
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get a recommendation: %s", err))
		return
	}

	respondWithJSON(w, 200, recommendation)
}

func (apiCfg apiConfig) handlerCreateRecommendations(w http.ResponseWriter, r *http.Request) {

}
