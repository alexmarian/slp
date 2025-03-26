package handlers

type ChirpValidationRequest struct {
	Body string `json:"body"`
}
type ChirpValidationResponse struct {
	CleanedBody string `json:"cleaned_body"`
	Error       string `json:"error"`
}
