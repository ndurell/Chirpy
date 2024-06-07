package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdVal := r.PathValue("chirpId")
	if chirpIdVal == "" {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	chirpId, err := strconv.Atoi(chirpIdVal)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.db.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	if chirp.Body == "" {
		respondWithError(w, 404, "Chirp not found")
		return
	}
	respondWithJSON(w, 200, chirp)
}
