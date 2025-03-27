package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleCreateUser(cfg *ApiConfig) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		request := CreateUserRequest{}
		err := decoder.Decode(&request)
		if err != nil {
			var errors = fmt.Sprintf("Error decoding create user request: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusBadRequest, errors)
			return
		}
		user, err := cfg.Db.CreateUser(req.Context(), request.Email)
		if err != nil {
			var errors = fmt.Sprintf("Error creating user: %s", err)
			log.Printf(errors)
			respondWithError(rw, http.StatusInternalServerError, errors)
			return
		}
		response := CreateUserResponse{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		respondWithJSON(rw, http.StatusCreated, response)
	}
}
