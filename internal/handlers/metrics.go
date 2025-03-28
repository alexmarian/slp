package handlers

import (
	"fmt"
	"github.com/alexmarian/slp/internal/auth"
	"github.com/alexmarian/slp/internal/database"
	"log"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
}

func (api *ApiConfig) IsDev() bool {
	return api.Platform == "dev"
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
func (cfg *ApiConfig) MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		userId, err := auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		next.ServeHTTP(w, addUserIdToContext(r, userId))
	}
}

func (cfg *ApiConfig) HandleMetrics(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	metrics := fmt.Sprintf(
		`<html>
					<body>
						<h1>Welcome, Chirpy Admin</h1>
						<p>Chirpy has been visited %d times!</p>
					</body>
				</html>`, cfg.fileserverHits.Load())
	rw.Write([]byte(metrics))
}

func (cfg *ApiConfig) HandleReset(rw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	if cfg.IsDev() {
		log.Printf("Metrics reset")
		cfg.Db.DeleteUsers(req.Context())
		rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Metrics reset"))
	} else {
		respondWithError(rw, http.StatusForbidden, "Reset only allowed in dev")
	}
}
