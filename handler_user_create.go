package main

import (
	"encoding/json"
	"net/http"
)

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	incomingUser := user{}
	err := decoder.Decode(&incomingUser)
	if err != nil {
		respondWithError(w, 400, "Invalid user")
		return
	}
	if len(incomingUser.Email) == 0 {
		respondWithError(w, 400, "Username is too short")
		return
	}
	user, err := cfg.db.CreateUser(incomingUser.Email, incomingUser.Password)
	if err != nil {
		respondWithError(w, 500, "Database error")
		return
	}
	respondWithJSON(w, 201, userResponse{
		Id:    user.Id,
		Email: user.Email,
	})
}
