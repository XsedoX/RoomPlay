package response

import (
	"encoding/json"
	"net/http"
)

const (
	encodingErrorMessage = "Could not encode response object."
)

type Success struct {
	Message *string `json:"message"  example:"null" extensions:"x-nullable"`
	Data    any     `json:"data" swaggertype:"object"`
}

func Ok(data any) Success {
	return Success{
		Message: nil,
		Data:    data,
	}
}

type Error struct {
	Message string  `json:"message"`
	Data    *string `json:"data" example:"null" extensions:"x-nullable"`
}

func Failure(message string) Error {
	return Error{
		Message: message,
		Data:    nil,
	}
}

func WriteJsonFailure(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Failure(message)); err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
	}
}

func WriteJsonSuccess[T any](w http.ResponseWriter, data T, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Ok(data)); err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
	}
}
