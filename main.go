package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	// Load the .env file into the environment variables
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to Open postgres: %v", err)
	}
	dbQueries := database.New(db)
	
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
	}
	
	mux := http.NewServeMux()	
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)
	
	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}


func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	type parameters struct {
		Body string `json:"body"`
    }

    type returnVals struct {
     	CleanedBody string `json:"cleaned_body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
	    log.Printf("Error decoding parameters: %s", err)
	    respondWithError(w, 500, "Something went wrong")
	    return
    }
    // check if we meet or exceed Chirpy char limit
    if len(params.Body) > 140 {
   		respondWithError(w, 400, "Chirp is too long")
     	return
    }
    cleanedChirp := cleanChirp(params.Body)
    respondWithJSON(w, 200, returnVals{CleanedBody: cleanedChirp})
    return
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	adminMessage := fmt.Sprintf("<html>\n  <body>\n   <h1>Welcome, Chirpy Admin</h1>\n   <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())
	w.Write([]byte(adminMessage))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
