package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform == "dev" {
		cfg.db.DeleteUsers(r.Context())
	} else {
		respondWithError(w, http.StatusForbidden, "403 Forbidden")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
	w.Write([]byte("Hits reset to 0"))
}