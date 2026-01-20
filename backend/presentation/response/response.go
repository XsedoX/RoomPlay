package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/customerrors"
	domainErrors "github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/infrastructure/validation"
	"github.com/go-playground/validator/v10"
)

const (
	encodingErrorMessage = "could not encode response object."
)

type Success struct {
	Data any `json:"data" swaggertype:"object" extensions:"x-nullable"`
}

type ProblemDetails struct {
	Type             string            `json:"type" example:"Error code unique for the error"`
	ValidationErrors map[string]string `json:"validationErrors" example:"{\"name\":\"too long\"}" swaggertype:"object" extensions:"x-nullable"`
	Title            string            `json:"title" example:"Short human readable description"`
	Description      string            `json:"description" example:"Longer human readable description"`
	Instance         string            `json:"instance" swaggertype:"string" example:"/api/v1/uri"`
	Status           int               `json:"status" swaggertype:"integer" example:"400"`
}

func WriteJsonFailure(w http.ResponseWriter, type1, title, description, instance string, statusCode int, errors ...map[string]string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.Header().Set("Content-Language", "en")

	var error1 map[string]string
	if len(errors) > 0 {
		error1 = errors[0]
	} else {
		error1 = nil
	}

	resp := &ProblemDetails{
		Type:             type1,
		ValidationErrors: error1,
		Title:            title,
		Description:      description,
		Instance:         instance,
		Status:           statusCode,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func WriteJsonApplicationFailure(w http.ResponseWriter, appErr error, instance string) {
	switch parsedErr := appErr.(type) {
	case *domainErrors.ValidationDomainError:
		WriteJsonFailure(w,
			parsedErr.Code,
			parsedErr.Title,
			parsedErr.Description,
			instance,
			http.StatusBadRequest)
		return
	case *customerrors.CustomError:
		WriteJsonFailure(w,
			parsedErr.Code,
			parsedErr.Title,
			parsedErr.Error(),
			instance,
			int(parsedErr.ErrorType))
		return
	default:
		WriteJsonFailure(w,
			"CustomError.CastingError",
			"Error is not applicationError",
			"Unexpected issue. Please try again.",
			instance,
			http.StatusInternalServerError)
		return
	}
}

func WriteJsonValidationFailure(w http.ResponseWriter, code, instance string, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		WriteJsonFailure(w,
			code,
			"Validation error occurred.",
			"One or more fields are not correctly filled.",
			instance,
			http.StatusBadRequest,
			validation.MapValidationErrors(validationErrs))
		return
	}
	WriteJsonFailure(w,
		"ValidationErrors.CastingError",
		"Error is not ValidationErrors",
		"Unexpected issue. Please try again.",
		instance,
		http.StatusInternalServerError)
}

func WriteJsonNoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func WriteJsonSuccess(w http.ResponseWriter, statusCode int, data ...any) {
	w.Header().Set("Content-Type", "application/json")

	resp := &Success{Data: nil}
	if len(data) > 0 {
		resp.Data = data[0]
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bytes)
}
