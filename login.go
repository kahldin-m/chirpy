package main

import (
	"encoding/json"
	"net/http"
	"time"

	"chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	expirationTime := time.Hour
	givenExpiration := time.Duration(params.ExpiresInSeconds) * time.Second
	if  givenExpiration > 0 && givenExpiration < 3600  {
		expirationTime = givenExpiration
	}

	pass, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !pass || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.serverSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error ocurred")
		return
	}


	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})
}
