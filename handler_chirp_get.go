package main

import (
	"database/sql"
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/jerslf/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	// Get query params
	authorIDStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Get all chirps first
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	// Filter by author_id (if provided)
	if authorIDStr != "" {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			http.Error(w, "invalid author_id", http.StatusBadRequest)
			return
		}

		filtered := []database.Chirp{}
		for _, c := range chirps {
			if c.UserID == authorID {
				filtered = append(filtered, c)
			}
		}
		chirps = filtered
	}

	// Sort by created_at (asc or desc)
	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	// Respond
	respondWithJSON(w, http.StatusOK, dbChirpsToAPI(chirps))
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, dbChirpToAPI(chirp))

}

func dbChirpToAPI(c database.Chirp) Chirp {
	return Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}

func dbChirpsToAPI(dbChirps []database.Chirp) []Chirp {
	result := make([]Chirp, len(dbChirps))
	for i, c := range dbChirps {
		result[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}
	return result
}
