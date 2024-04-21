package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := clearProfane(params.Body)
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

// Clean This Up Later
func clearProfane(s string) string {
	cleanedWords := []string{}
	words := strings.Split(s, " ")
	for _, word := range words {
		lowered_word := strings.ToLower(word)
		if lowered_word == "kerfuffle" || lowered_word == "sharbert" || lowered_word == "fornax" {
			word = "****"
		}
		cleanedWords = append(cleanedWords, word)
	}
	return strings.Join(cleanedWords, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
