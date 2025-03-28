package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alexmarian/slp/internal/auth"
	"github.com/alexmarian/slp/internal/database"
	"log"
	"net/http"
	"time"
)

func HandleLogin(cfg *ApiConfig) http.HandlerFunc {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
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
		seconds := 3600
		if request.ExpiresInSeconds != 0 {
			seconds = request.ExpiresInSeconds
		}
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Error creating refresh token")
			return
		}
		rt, err := cfg.Db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: user.ID,
		})
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Error creating refresh token")
		}
		token, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(seconds)*time.Second)
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
			Token:        token,
			RefreshToken: rt.Token,
		}
		respondWithJSON(rw, http.StatusOK, usr)
	}
}

func HandleRefresh(cfg *ApiConfig) http.HandlerFunc {
	type response struct {
		Token string `json:"token,omitempty"`
	}
	return func(rw http.ResponseWriter, req *http.Request) {
		refreshToken, err := auth.GetBearerToken(req.Header)
		if err != nil {
			respondWithError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}
		rt, err := cfg.Db.GetValidRefreshToken(req.Context(), refreshToken)
		if err != nil {
			respondWithError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}
		token, err := auth.MakeJWT(rt.UserID, cfg.Secret, 3600*time.Second)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Error creating token")
			return
		}
		resp := response{
			Token: token,
		}
		respondWithJSON(rw, http.StatusOK, resp)
	}
}
