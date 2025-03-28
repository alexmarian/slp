package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alexmarian/slp/internal/auth"
	"log"
	"net/http"
	"time"
)

func HandleLogin(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token,omitempty"`
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		request := parameters{}
		err := decoder.Decode(&request)
		if err != nil {
			var errors = fmt.Sprintf("Error decoding login user request: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusBadRequest, errors)
			return
		}
		user, err := cfg.Db.GetUserByEmail(req.Context(), request.Email)
		if err != nil {
			var errors = fmt.Sprintf("Error getting user: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusInternalServerError, errors)
			return
		}
		err = auth.CheckPasswordHash(request.Password, user.HashedPassword)
		if err != nil {
			respondWithError(rw, http.StatusUnauthorized, "Incorrect email or password")
			return
		}
		token, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(request.ExpiresInSeconds)*time.Second)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Error creating token")
			return
		}
		usr := response{
			User: User{
				Id:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
			Token: token,
		}
		respondWithJSON(rw, http.StatusOK, usr)
	}
}
