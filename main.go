package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := &apiConfig{fileserverHits: 0}

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fileHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileHandler)))
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("/api/reset", apiCfg.handleMetricsReset)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type: ", "text/html")
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`
	<html>

	<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
	</body>

	</html>`, cfg.fileserverHits)
	w.Write([]byte(html))
}
