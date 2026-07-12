package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/domain/domain_errors"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	encodingErrorMessage = "could not encode response object."
)

type Success struct {
	Data any                        `json:"data" swaggertype:"object" extensions:"x-nullable"`
	Meta *page_meta_dto.PageMetaDto `json:"meta" swaggertype:"object" extensions:"x-nullable"`
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

func WriteJsonDecodingFailure(w http.ResponseWriter, code string, bodyDecodeErr error, instance string) {
	WriteJsonFailure(w,
		code,
		"Unexpected issue occurred while decoding the request body.",
		"Something went wrong internally. Errror code: "+code,
		instance,
		http.StatusBadRequest,
	)
}

func WriteJsonApplicationFailure(w http.ResponseWriter, appErr error, instance string) {
	if vErr, ok := errors.AsType[domain_errors.DomainError](appErr); ok {
		WriteJsonFailure(w,
			vErr.Code,
			"Some error occurred.",
			vErr.Description,
			instance,
			http.StatusBadRequest,
		)
		return
	}
	if cErr, ok := errors.AsType[*application_error.ApplicationError](appErr); ok {
		WriteJsonFailure(w,
			cErr.Code,
			cErr.Title,
			cErr.Title,
			instance,
			int(cErr.ErrorType),
		)
		return
	}
	WriteJsonFailure(w,
		"CustomError.CastingError",
		"Error is not applicationError",
		"Unexpected issue. Please try again.",
		instance,
		http.StatusInternalServerError,
	)
}

func WriteJsonValidationFailure(w http.ResponseWriter, code, instance string, err error) {
	if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
		WriteJsonFailure(w,
			code,
			"Validation error occurred.",
			"One or more fields are not correctly filled",
			instance,
			http.StatusUnprocessableEntity,
			setup_validation.MapValidationErrors(validationErrs))
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

func WriteJsonCreated(w http.ResponseWriter, id uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")

	resp := &Success{Data: id}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

func WriteJsonSuccess(w http.ResponseWriter, data any, meta ...page_meta_dto.PageMetaDto) {
	w.Header().Set("Content-Type", "application/json")

	resp := &Success{Data: data}
	if len(meta) > 0 {
		resp.Meta = &meta[0]
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, encodingErrorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
