package util

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func RespondWithSuccess(w http.ResponseWriter, status int, message string, data any) {
	WriteJSON(w, status, SuccessResponse{Message: message, Data: data})
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Message: message})
}
