package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Gfarf/Chirpy/internal/auth"
	"github.com/Gfarf/Chirpy/internal/database"
)

func UserMapping(user *database.User) User {
	res := User{}
	res.ID = user.ID
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
