package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Gfarf/Chirpy/internal/auth"
	"github.com/Gfarf/Chirpy/internal/database"
	"github.com/google/uuid"
)

func UserMapping(user *database.User) User {
	res := User{}
	res.ID = uuid.UUID(user.ID)
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
	res.Email = user.Email
	return res
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type receiveUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	userEmail := receiveUser{}
	err := decoder.Decode(&userEmail)
	if err != nil {
		log.Printf("Error decoding user e-mail: %s", err)
		w.WriteHeader(500)
		return
	}
	pWord, err := auth.HashPassword(userEmail.Password)
	if err != nil {
		log.Printf("Error hashing user password: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: userEmail.Email, HashedPassword: pWord})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	res := UserMapping(&user)
	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	type receiveUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	userEmail := receiveUser{}
	err := decoder.Decode(&userEmail)
	if err != nil {
		log.Printf("Error decoding user e-mail: %s", err)
		w.WriteHeader(500)
		return
	}
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Bearer Token: %s", err)
		w.WriteHeader(401)
		return
	}
	UserID, err := auth.ValidateJWT(jwt, cfg.secretString)
	if err != nil {
		log.Printf("Invalid JWT token in chirping: %s", err)
		w.WriteHeader(401)
		return
	}
	pWord, err := auth.HashPassword(userEmail.Password)
	if err != nil {
		log.Printf("Error hashing user password: %s", err)
		w.WriteHeader(500)
		return
	}
	err = cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{Email: userEmail.Email, HashedPassword: pWord, ID: UserID})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	u, err := cfg.dbQueries.GetUserByEmail(r.Context(), userEmail.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(500)
		return
	}
	res := UserMapping(&u)
	res.LoginToken, err = auth.MakeJWT(res.ID, cfg.secretString, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(500)
		return
	}
	res.RefreshToken, err = cfg.SaveRefreshToken(res.ID)
	if err != nil {
		log.Printf("Error creating Refresh Token: %s", err)
		w.WriteHeader(500)
		return
	}
	dat, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
