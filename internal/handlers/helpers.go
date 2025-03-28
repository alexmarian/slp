package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
)

const userContextKey = "userID"

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func addUserIdToContext(req *http.Request, userID uuid.UUID) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), userContextKey, userID))
}

func getUserIdFromContext(req *http.Request) uuid.UUID {
	return req.Context().Value(userContextKey).(uuid.UUID)
}
