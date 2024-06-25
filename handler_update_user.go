package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type updateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	incomingUser := updateUser{}
	err := decoder.Decode(&incomingUser)
	if err != nil {
		respondWithError(w, 400, "Invalid user")
		return
	}
	bearerToken := r.Header.Get("Authorization")
	bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")
	claims := jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "Chirpy",
		Subject:   strconv.Itoa(1),
	}
	res, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		fmt.Println("Error parsing token", err, bearerToken)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userIdStr, err := res.Claims.GetSubject()
	if err != nil {
		fmt.Println("Error getting subject", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println("Error converting user id", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	user := cfg.db.GetUserById(userId)
	if user == nil {
		respondWithError(w, 404, "User not found")
		return
	}
	err = cfg.db.UpdateUser(user.Id, incomingUser.Email, incomingUser.Password)
	if err != nil {
		respondWithError(w, 500, "Error Updating user")
		return
	}

	respondWithJSON(w, 200, userResponse{
		Id:    user.Id,
		Email: incomingUser.Email,
	})

}
