package webhooks

import (
	"encoding/json"
	"fmt"
	"github.com/alexmarian/slp/internal/auth"
	"github.com/alexmarian/slp/internal/handlers"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func HandleUpdateChirpyRed(cfg *handlers.ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := auth.GetApiKey(r.Header)
		if err != nil || key != cfg.PolkaKey {
			w.WriteHeader(http.StatusUnauthorized)
			return

		}
		type params struct {
			Event string `json:"event"`
			Data  struct {
				UserId uuid.UUID `json:"user_id"`
			} `json:"data"`
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		reqParam := params{}
		err = decoder.Decode(&reqParam)
		if err != nil {
			var errors = fmt.Sprintf("Error decoding chirp: %s", err)
			log.Printf(errors)
			handlers.RespondWithError(w, http.StatusBadRequest, errors)
			return
		}
		if reqParam.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		err = cfg.Db.UpdateUserToChirpyRed(r.Context(), reqParam.Data.UserId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	}
}
