package handlers

import (
	"github.com/alexmarian/slp/internal/auth"
	"net/http"
)

func HandleRevokeRefreshToken(cfg *ApiConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		refreshToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			RespondWithError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}
		err = cfg.Db.RevokeRefreshToken(req.Context(), refreshToken)
		if err != nil {
			RespondWithError(rw, http.StatusInternalServerError, "Error revoking token")
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
