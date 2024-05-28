package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

type validChirp struct {
	CleanBody string `json:"cleaned_body"`
}

var BAD_WORDS = []string{"kerfuffle", "sharbert", "fornax"}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	incomingChirp := chirp{}
	err := decoder.Decode(&incomingChirp)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp")
		return
	}
	if len(incomingChirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	respBody := validChirp{
		CleanBody: cleanChirp(incomingChirp),
	}
	respondWithJSON(w, 200, respBody)
}

func cleanChirp(c chirp) string {

	words := strings.Split(c.Body, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if slices.Contains(BAD_WORDS, lowered) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
