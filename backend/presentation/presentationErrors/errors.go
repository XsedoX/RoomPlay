package presentationErrors

import (
	"encoding/json"
	"net/http"

	"xsedox.com/main/application/response"
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
