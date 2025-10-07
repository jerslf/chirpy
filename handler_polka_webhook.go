package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	// Ignore events other than "user.upgraded"
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Update user as is_chirpy_red
	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user", err)
		return
	}

	// Everything OK
	w.WriteHeader(http.StatusNoContent)
}
