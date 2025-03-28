package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alexmarian/slp/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sort"
	"strings"
)

func HandleCreateChirp(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		chirp := ChirpCreationRequest{}
		err := decoder.Decode(&chirp)
		if err != nil {
			var errors = fmt.Sprintf("Error decoding chirp: %s", err)
			log.Printf(errors)
			RespondWithError(w, http.StatusBadRequest, errors)
			return
		}
		if len(chirp.Body) > 140 {
			response := Chirp{}
			response.Body = "Chirp is too long"
			RespondWithJSON(w, http.StatusBadRequest, response)
			return
		}
		badWords := map[string]struct{}{
			"kerfuffle": {},
			"sharbert":  {},
			"fornax":    {},
		}
		ccp := database.CreateChirpParams{
			Body:   chirp.Body,
			UserID: GetUserIdFromContext(r),
		}
		createChirp, err := cfg.Db.CreateChirp(r.Context(), ccp)
		if err != nil {
			var errors = fmt.Sprintf("Error creating chirp: %s", err)
			log.Printf(errors)
			RespondWithError(w, http.StatusInternalServerError, errors)
			return
		}
		response := Chirp{
			Id:        createChirp.ID,
			CreatedAt: createChirp.CreatedAt,
			UpdatedAt: createChirp.UpdatedAt,
			Body:      createChirp.Body,
			UserId:    createChirp.UserID,
		}
		response.Body = cleanBody(chirp.Body, badWords)
		RespondWithJSON(w, http.StatusCreated, response)
	}
}

func cleanBody(body string, badWords map[string]struct{}) string {
	tokens := strings.Split(body, " ")
	for i, token := range tokens {
		if _, ok := badWords[strings.ToLower(token)]; ok {
			tokens[i] = "****"
		}
	}
	return strings.Join(tokens, " ")
}

func HandleGetChirps(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authorID := r.URL.Query().Get("author_id")
		var chirps []database.Chirp
		if authorID != "" {
			authorUuid, err := uuid.Parse(authorID)
			if err != nil {
				var errors = fmt.Sprintf("Error parsing author ID: %s", err)
				log.Printf(errors)
				RespondWithError(w, http.StatusBadRequest, errors)
				return
			}
			chirps, err = cfg.Db.GetAllChirpsByAuthorId(r.Context(), authorUuid)
			if err != nil {
				var errors = fmt.Sprintf("Error getting chirps: %s", err)
				log.Printf(errors)
				RespondWithError(w, http.StatusInternalServerError, errors)
				return
			}
		} else {
			var err error
			chirps, err = cfg.Db.GetAllChirps(r.Context())
			if err != nil {
				var errors = fmt.Sprintf("Error getting chirps: %s", err)
				log.Printf(errors)
				RespondWithError(w, http.StatusInternalServerError, errors)
				return
			}
		}
		var response []Chirp
		sortDirection := r.URL.Query().Get("sort")
		if sortDirection == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		}
		for _, chirp := range chirps {
			response = append(response, Chirp{
				Id:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			})
		}
		RespondWithJSON(w, http.StatusOK, response)
	}
}
func HandleGetChirp(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpID, parseIdErr := uuid.Parse(r.PathValue("chirpID"))
		if parseIdErr != nil {
			var errors = fmt.Sprintf("Error parsing chirp ID: %s", parseIdErr)
			log.Printf(errors)
			RespondWithError(w, http.StatusBadRequest, errors)
			return
		}
		chirp, err := cfg.Db.GetChirpById(r.Context(), chirpID)
		if err != nil {
			log.Printf("Error getting chirp: %s", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		RespondWithJSON(w, http.StatusOK, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
}
func HandleDeleteChirp(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpID, parseIdErr := uuid.Parse(r.PathValue("chirpID"))
		if parseIdErr != nil {
			var errors = fmt.Sprintf("Error parsing chirp ID: %s", parseIdErr)
			log.Printf(errors)
			RespondWithError(w, http.StatusBadRequest, errors)
			return
		}
		chirp, err := cfg.Db.GetChirpById(r.Context(), chirpID)
		if err != nil {
			log.Printf("Error getting chirp: %s", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		userId := GetUserIdFromContext(r)
		if chirp.UserID != userId {
			log.Printf("Error deleting chirp: %s", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		err = cfg.Db.DeleteChirpById(r.Context(), chirp.ID)
		if err != nil {
			var errors = fmt.Sprintf("Error deleting chirp: %s", err)
			log.Printf(errors)
			RespondWithError(w, http.StatusInternalServerError, errors)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
