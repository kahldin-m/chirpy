package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"

	"chirpy/internal/database"
)

type returnVals struct {
    	ID uuid.UUID `json:"id"`
     	CreatedAt time.Time `json:"created_at"`
      	UpdatedAt time.Time `json:"updated_at"`
       	Body string `json:"body"`
        UserID uuid.UUID `json:"user_id"`
    }

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
    }

    // decoding both fields in parameters
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
	    log.Printf("Error decoding parameters: %s", err)
	    respondWithError(w, 500, "Something went wrong")
	    return
    }
    // check if we meet or exceed Chirpy char limit (validate)
    if len(params.Body) > 140 {
   		respondWithError(w, 400, "Chirp is too long")
     	return
    }
    cleanedChirp := cleanChirp(params.Body)
    finalChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
    	Body: cleanedChirp,
     	UserID: params.UserID,
    })
    if err != nil {
    	respondWithError(w, 400, "Error creating chirp")
     	return
    }
    respondWithJSON(w, 201, returnVals{
    	ID: finalChirp.ID,
     	CreatedAt: finalChirp.CreatedAt,
      	UpdatedAt: finalChirp.UpdatedAt,
      	Body: finalChirp.Body,
      	UserID: finalChirp.UserID,
    })
    return
}


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong retrieving chirps")
		return
	}
	results := []returnVals{}
	for _, chirp := range(chirps) {
		results = append(results, returnVals{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}
	respondWithJSON(w, 200, results)
	return
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Invalid ID or ID not found")
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "Chirp with that ID not found")
		return
	}
	result := returnVals{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, 200, result)
	return	
}
