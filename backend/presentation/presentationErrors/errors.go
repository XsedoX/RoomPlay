package presentationErrors

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/presentation/response"
)

const (
	encodingErrorMessage = "Could not encode response object."
)

func WriteJsonFailure(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response.Failure(message)); err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
	}
}

func WriteJsonSuccess[T any](w http.ResponseWriter, data T, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response.Ok(data)); err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
	}
}
