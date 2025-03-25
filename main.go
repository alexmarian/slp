package main

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
	mux.HandleFunc("/healthz", handlers.HandleHealthz)
	mux.HandleFunc("/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("/reset", apiCfg.HandleReset)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
