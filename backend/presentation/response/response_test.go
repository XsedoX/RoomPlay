package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/application/dtos/page_meta_dto"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup_validation.Initialize()
	m.Run()
}

func TestWriteJsonFailure(t *testing.T) {
	t.Run("Basic Failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteJsonFailure(w, "type", "title", "description", "instance", http.StatusBadRequest)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "type", response.Type)
		assert.Equal(t, "title", response.Title)
		assert.Equal(t, "description", response.Description)
		assert.Equal(t, "instance", response.Instance)
		assert.Equal(t, http.StatusBadRequest, response.Status)
		assert.Nil(t, response.ValidationErrors)
	})

	t.Run("Failure with Validation Errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		validationErrors := map[string]string{"field": "error"}
		WriteJsonFailure(w, "type", "title", "description", "instance", http.StatusBadRequest, validationErrors)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, validationErrors, response.ValidationErrors)
		assert.Equal(t, http.StatusBadRequest, response.Status)
	})

	t.Run("Failure with Multiple Error Maps", func(t *testing.T) {
		w := httptest.NewRecorder()
		validationErrors1 := map[string]string{"field1": "error1"}
		validationErrors2 := map[string]string{"field2": "error2"}
		WriteJsonFailure(w, "type", "title", "description", "instance", http.StatusBadRequest, validationErrors1, validationErrors2)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)
		assert.Equal(t, validationErrors1, response.ValidationErrors)
	})
}

func TestWriteJsonApplicationFailure(t *testing.T) {
	t.Run("Valid Custom Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		customErr := application_error.NewApplicationError("code", "title", errors.New("inner error"), application_error_type.Validation)
		WriteJsonApplicationFailure(w, customErr, "instance")

		assert.Equal(t, http.StatusBadRequest, w.Code) // Validation type maps to 400

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "code", response.Type)
		assert.Equal(t, "title", response.Title)
		assert.Contains(t, response.Description, response.Title)
		assert.Equal(t, "instance", response.Instance)
	})

	t.Run("Non-Custom Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		stdErr := errors.New("standard error")
		WriteJsonApplicationFailure(w, stdErr, "instance")

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "CustomError.CastingError", response.Type)
		assert.Equal(t, "Error is not applicationError", response.Title)
	})
}

func TestWriteJsonValidationFailure(t *testing.T) {
	// Helper struct to generate validation errors
	type TestStruct struct {
		Name string `validate:"required"`
	}

	t.Run("Valid Validation Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		validate := validator.New()
		err := validate.Struct(TestStruct{}) // Should fail
		assert.Error(t, err)

		WriteJsonValidationFailure(w, "code", "instance", err)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var response ProblemDetails
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "code", response.Type)
		assert.Equal(t, "Validation error occurred.", response.Title)
		assert.NotEmpty(t, response.ValidationErrors)
		assert.Contains(t, response.ValidationErrors, "Name")
	})

	t.Run("Non-Validation Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		stdErr := errors.New("standard error")
		WriteJsonValidationFailure(w, "code", "instance", stdErr)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ValidationErrors.CastingError", response.Type)
		assert.Equal(t, "Error is not ValidationErrors", response.Title)
	})
}

func TestWriteJsonNoContent(t *testing.T) {
	t.Run("Standard Execution", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteJsonNoContent(w)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Empty(t, w.Body.String())
	})
}

func TestWriteJsonSuccess(t *testing.T) {
	t.Run("Success with meta", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}
		meta := page_meta_dto.PageMetaDto{
			NextPageToken:     new(gofakeit.ID()),
			PreviousPageToken: new(gofakeit.ID()),
			PageSize:          10,
			HasNextPage:       false,
		}
		WriteJsonSuccess(w, data, meta)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		responseMeta := response.Meta
		assert.Equal(t, &meta, responseMeta)
		// Need to cast response.Data to map to check content
		responseData, ok := response.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "value", responseData["key"])
	})
	t.Run("Success without meta", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}
		WriteJsonSuccess(w, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Need to cast response.Data to map to check content
		responseData, ok := response.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "value", responseData["key"])
		assert.Nil(t, response.Meta)
	})

	t.Run("Success StatusCreated", func(t *testing.T) {
		w := httptest.NewRecorder()
		uuidToReturn := uuid.New()
		WriteJsonCreated(w, uuidToReturn)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uuidToReturn.String(), response.Data)
	})
}
