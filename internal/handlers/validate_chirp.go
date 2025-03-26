package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HandleValidateChirp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := ChirpValidationRequest{}
	err := decoder.Decode(&chirp)
	if err != nil {
		var errors = fmt.Sprintf("Error decoding chirp: %s", err)
		log.Printf(errors)
		respondWithError(rw, http.StatusBadRequest, errors)
		return
	}
	response := ChirpValidationResponse{}
	if len(chirp.Body) > 140 {
		response.Error = "Chirp is too long"
		respondWithJSON(rw, http.StatusBadRequest, response)
		return
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	response.CleanedBody = cleanBody(chirp.Body, badWords)
	respondWithJSON(rw, http.StatusOK, response)
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
