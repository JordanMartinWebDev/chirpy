package main

import (
	"errors"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func ValidateChirp(s string) (string, error) {
	const maxChirpLength = 140

	if len(s) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := clearProfane(s, badWords)
	return cleaned, nil
}

func clearProfane(s string, badWords map[string]struct{}) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
