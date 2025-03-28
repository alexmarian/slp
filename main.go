package main

import (
	"database/sql"
	"github.com/alexmarian/slp/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)
import (
	"github.com/alexmarian/slp/internal/handlers"
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load(".env")
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	dbQueries := database.New(db)
	apiCfg := &handlers.ApiConfig{
		Db:       dbQueries,
		Platform: platform,
		Secret:   secret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlers.HandleHealthz)

	mux.HandleFunc("POST /api/chirps", handlers.HandleCreateChirp(apiCfg))
	mux.HandleFunc("GET /api/chirps", handlers.HandleGetChirps(apiCfg))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.HandleGetChirp(apiCfg))
	mux.HandleFunc("POST /api/users", handlers.HandleCreateUser(apiCfg))
	mux.HandleFunc("POST /api/login", handlers.HandleLogin(apiCfg))

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
