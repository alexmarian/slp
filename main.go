package main

import _ "github.com/lib/pq"
import (
	"github.com/alexmarian/slp/internal/handlers"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &handlers.ApiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlers.HandleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandleValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
