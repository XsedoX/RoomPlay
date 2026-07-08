package response

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XsedoX/RoomPlay/application/application_error"
	"github.com/XsedoX/RoomPlay/application/application_error/application_error_type"
	"github.com/XsedoX/RoomPlay/presentation/setup_validation"
	"github.com/go-playground/validator/v10"
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
		assert.Contains(t, response.Description, "inner error")
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

		assert.Equal(t, http.StatusBadRequest, w.Code)

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
	t.Run("Success with Data", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}
		WriteJsonSuccess(w, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Need to cast response.Data to map to check content
		responseData, ok := response.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "value", responseData["key"])
	})

	t.Run("Success with Nil Data", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteJsonSuccess(w, http.StatusOK, nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Nil(t, response.Data)
	})

	t.Run("Success with No Data Arguments", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteJsonSuccess(w, http.StatusOK)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Nil(t, response.Data)
	})

	t.Run("Success with Multiple Data Arguments", func(t *testing.T) {
		w := httptest.NewRecorder()
		data1 := "data1"
		data2 := "data2"
		WriteJsonSuccess(w, http.StatusOK, data1, data2)

		var response Success
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, data1, response.Data)
	})

	t.Run("JSON Encoding Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		// math.Inf(1) is not valid JSON
		WriteJsonSuccess(w, http.StatusOK, math.Inf(1))

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "could not encode response object.")
	})
}
