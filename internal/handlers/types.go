package handlers

import (
	"github.com/google/uuid"
	"time"
)

type ChirpValidationRequest struct {
	Body string `json:"body"`
}
type ChirpValidationResponse struct {
	CleanedBody string `json:"cleaned_body"`
	Error       string `json:"error"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
}

type CreateUserResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
