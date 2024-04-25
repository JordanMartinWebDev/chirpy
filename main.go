package main

import (
	"log"
	"net/http"

	"github.com/jordanmartinwebdev/chirpy/internal/database"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))

	//API Handlers
	mux.HandleFunc("GET /api/reset", apiCfg.handleMetricsReset)
	mux.HandleFunc("GET /api/healthz", apiCfg.handleReadiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGet)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}
