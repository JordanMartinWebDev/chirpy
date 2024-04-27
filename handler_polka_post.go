package main

import (
	"encoding/json"
	"net/http"

	"github.com/jordanmartinwebdev/chirpy/internal/auth"
)

type Data struct {
	UserID int `json:"user_id"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if key != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	err = cfg.DB.EnableChirpyRed(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find user to upgrade")
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
