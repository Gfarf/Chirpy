package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Gfarf/Chirpy/internal/auth"
	"github.com/Gfarf/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) SaveRefreshToken(userId uuid.UUID) (string, error) {
	newToken, _ := auth.MakeRefreshToken()
	rtk, err := cfg.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{ID: newToken, UserID: userId})
	if err != nil {
		return "", err
	}
	return rtk.ID, nil
}

func (cfg *apiConfig) ValidateRefreshToken(token string) (RefreshToken, error) {
	rToken, err := cfg.dbQueries.ValidateToken(context.Background(), token)
	res := RefreshToken{}
	if err != nil {
		return res, err
	}
	res.ID = rToken.ID
	res.CreatedAt = rToken.CreatedAt
	res.UpdatedAt = rToken.UpdatedAt
	res.UserID = rToken.UserID
	res.ExpiresAt = rToken.ExpiresAt
	if rToken.RevokedAt.Valid {
		res.RevokedAt = rToken.RevokedAt.Time
	}
	res.Revoked = rToken.RevokedAt.Valid
	return res, nil
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Refresh Token: %s", err)
		w.WriteHeader(500)
		return
	}
	refreshToken, err := cfg.ValidateRefreshToken(token)
	if err != nil {
		log.Printf("Invalid refresh token on handler refresh: %s", err)
		w.WriteHeader(401)
		return
	}
	if refreshToken.Revoked {
		log.Printf("Revoked token: %s", err)
		w.WriteHeader(401)
		return
	}

	res, err := auth.MakeJWT(refreshToken.UserID, cfg.secretString, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		w.WriteHeader(500)
		return
	}

	fRes2 := struct {
		Token string `json:"token"`
	}{
		Token: res,
	}
	dat, err := json.Marshal(fRes2)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Refresh Token: %s", err)
		w.WriteHeader(500)
		return
	}
	_, err = cfg.ValidateRefreshToken(token)
	if err != nil {
		log.Printf("Invalid refresh token on revoking: %s", err)
		w.WriteHeader(401)
		return
	}
	err = cfg.dbQueries.RevokeToken(r.Context(), token)
	if err != nil {
		log.Printf("Error revoking token: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
