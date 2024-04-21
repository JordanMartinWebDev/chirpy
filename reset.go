package main

import (
	"net/http"
)

func (cfg *apiConfig) handleMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.Header().Add("Content-Type: ", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	reset := "Metrics Have Been Reset"
	w.Write([]byte(reset))
}
