package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userLogin struct {
	ExpiresInSeconds *int `json:"expires_in_seconds"`
	user
}

type userResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
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

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	incomingUser := userLogin{}
	err := decoder.Decode(&incomingUser)
	if err != nil {
		respondWithError(w, 400, "Invalid user")
		return
	}
	if len(incomingUser.Email) == 0 {
		respondWithError(w, 400, "Username is too short")
		return
	}
	user, err := cfg.db.GetUser(incomingUser.Email)
	if err != nil {
		fmt.Println("Error getting user", err)
		respondWithError(w, 500, "Database error")
		return
	}
	if user == nil {
		respondWithError(w, 404, "User not found")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(incomingUser.Password))
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	expires := 24 * 60 * 60
	if incomingUser.ExpiresInSeconds != nil && *incomingUser.ExpiresInSeconds < expires {
		expires = *incomingUser.ExpiresInSeconds
	}
	claims := jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires) * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "Chirpy",
		Subject:   strconv.Itoa(user.Id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		fmt.Println("Error signing token", err)
		respondWithError(w, 500, "Database error")
		return
	}

	respondWithJSON(w, 200, userResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: tokenString,
	})
}
