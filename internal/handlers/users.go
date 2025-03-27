package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/alexmarian/slp/internal/auth"
	"github.com/alexmarian/slp/internal/database"
	"log"
	"net/http"
)

func HandleCreateUser(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		request := UserRequest{}
		err := decoder.Decode(&request)
		if err != nil {
			var errors = fmt.Sprintf("Error decoding create user request: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusBadRequest, errors)
			return
		}
		password, err := auth.HashPassword(request.Password)
		if err != nil {
			var errors = fmt.Sprintf("Error hashing password: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusInternalServerError, errors)
			return
		}
		user, err := cfg.Db.CreateUser(req.Context(), database.CreateUserParams{request.Email, password})
		if err != nil {
			var errors = fmt.Sprintf("Error creating user: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusInternalServerError, errors)
			return
		}
		response := User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		respondWithJSON(rw, http.StatusCreated, response)
	}
}
