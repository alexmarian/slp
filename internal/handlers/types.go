package handlers

import (
	"github.com/google/uuid"
	"time"
)

type ChirpCreationRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}
type ChirpCreationResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    string    `json:"user_id"`
	Error     string    `json:"error"`
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
