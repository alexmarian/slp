package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleValidateChirp(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirp := ChirpValidationRequest{}
	response := ChirpValidationResponse{}
	err := decoder.Decode(&chirp)
	if err != nil {
		response.Error = fmt.Sprintf("Error decoding chirp: %s", err)
		log.Printf(response.Error)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		if len(chirp.Body) > 140 {
			response.Valid = false
			response.Error = "Chirp is too long"
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			response.Valid = true
			rw.WriteHeader(http.StatusOK)
		}
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	rw.Write(data)
}
