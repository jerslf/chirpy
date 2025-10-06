package main

import (
	"net/http"
	"time"

	"github.com/jerslf/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}

	// Extract the token from the Authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization header", err)
		return
	}

	dbToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil || dbToken.RevokedAt.Valid || time.Now().After(dbToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	// Create new access token for the user
	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Token: accessToken})

}
