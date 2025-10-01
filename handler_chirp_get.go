package main

import (
	"net/http"

	"github.com/jerslf/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	respondWithJSON(w, http.StatusOK, dbChirpsToAPI(chirps))
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
