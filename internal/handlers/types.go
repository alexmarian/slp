package handlers

type ChirpValidationRequest struct {
	Body string `json:"body"`
}
type ChirpValidationResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}
