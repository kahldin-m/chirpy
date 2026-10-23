package main

import (
	"encoding/json"
	"net/http"

	"chirpy/internal/auth"
	"chirpy/internal/database"
)


type parameters struct {
	Password   string `json:"password"`
	Email      string `json:"email"`
}


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	pword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Hashing failed. Blame the salt mines")
		return
	}
	
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: pword,
		Email: params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
