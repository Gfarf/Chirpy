package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Gfarf/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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
	hashedP, err := cfg.dbQueries.GetHashedPasswordByEmail(r.Context(), userEmail.Email)
	if err != nil {
		log.Printf("Incorrect email or password: %s", err)
		w.WriteHeader(401)
		return
	}
	err = auth.CheckPasswordHash(userEmail.Password, hashedP)
	if err != nil {
		log.Printf("Incorrect email or password: %s", err)
		w.WriteHeader(401)
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
