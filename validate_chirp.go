package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/ndurell/Chirpy/internal/database"
)

type chirp struct {
	Body string `json:"body"`
}

type validChirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

var BAD_WORDS = []string{"kerfuffle", "sharbert", "fornax"}

var chirpCount = 0

func createChirp(w http.ResponseWriter, r *http.Request) {

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
		Id:   chirpCount,
		Body: cleanChirp(incomingChirp),
	}
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	log.Printf("Creating chirp: %s", respBody.Body)
	chirp, err := db.CreateChirp(respBody.Body)
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	log.Printf("Created chirp: %s", chirp.Body)

	respondWithJSON(w, 201, chirp)
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

func getChirps(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	respondWithJSON(w, 200, chirps)
}
