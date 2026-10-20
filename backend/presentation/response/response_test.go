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
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setup_validation.Initialize()
	m.Run()
}

func TestWriteJsonFailure(t *testing.T) {
	t.Run("Basic Failure", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteJsonFailure(w, "type", "title", "description", "instance", http.StatusBadRequest)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "type", response.Type)
		require.Equal(t, "title", response.Title)
		require.Equal(t, "description", response.Description)
		require.Equal(t, "instance", response.Instance)
		require.Equal(t, http.StatusBadRequest, response.Status)
		require.Nil(t, response.ValidationErrors)
	})

	t.Run("Failure with Validation Errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		validationErrors := map[string]string{"field": "error"}
		WriteJsonFailure(w, "type", "title", "description", "instance", http.StatusBadRequest, validationErrors)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, validationErrors, response.ValidationErrors)
		require.Equal(t, http.StatusBadRequest, response.Status)
	})

	t.Run("Failure with Multiple Error Maps", func(t *testing.T) {
		w := httptest.NewRecorder()
		validationErrors1 := map[string]string{"field1": "error1"}
		validationErrors2 := map[string]string{"field2": "error2"}
		WriteJsonFailure(w, "type", "title", "description", "instance", http.StatusBadRequest, validationErrors1, validationErrors2)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, response.Status)
		require.Equal(t, validationErrors1, response.ValidationErrors)
	})
}

func TestWriteJsonApplicationFailure(t *testing.T) {
	t.Run("Valid Custom Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		customErr := application_error.NewApplicationError("code", "title", errors.New("inner error"), application_error_type.Validation)
		WriteJsonApplicationFailure(w, customErr, "instance")

		require.Equal(t, http.StatusBadRequest, w.Code) // Validation type maps to 400

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "code", response.Type)
		require.Equal(t, "title", response.Title)
		require.Contains(t, response.Description, response.Title)
		require.Equal(t, "instance", response.Instance)
	})

	t.Run("Non-Custom Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		stdErr := errors.New("standard error")
		WriteJsonApplicationFailure(w, stdErr, "instance")

		require.Equal(t, http.StatusInternalServerError, w.Code)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "CustomError.CastingError", response.Type)
		require.Equal(t, "Error is not applicationError", response.Title)
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
		require.Error(t, err)

		WriteJsonValidationFailure(w, "code", "instance", err)

		require.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var response ProblemDetails
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "code", response.Type)
		require.Equal(t, "Validation error occurred.", response.Title)
		require.NotEmpty(t, response.ValidationErrors)
		require.Contains(t, response.ValidationErrors, "Name")
	})

	t.Run("Non-Validation Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		stdErr := errors.New("standard error")
		WriteJsonValidationFailure(w, "code", "instance", stdErr)

		require.Equal(t, http.StatusInternalServerError, w.Code)

		var response ProblemDetails
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, "ValidationErrors.CastingError", response.Type)
		require.Equal(t, "Error is not ValidationErrors", response.Title)
	})
}

func TestWriteJsonNoContent(t *testing.T) {
	t.Run("Standard Execution", func(t *testing.T) {
		w := httptest.NewRecorder()
		WriteJsonNoContent(w)

		require.Equal(t, http.StatusNoContent, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.Empty(t, w.Body.String())
	})
}

func TestWriteJsonSuccess(t *testing.T) {
	t.Run("Success with meta", func(t *testing.T) {
		w := httptest.NewRecorder()
		type testResult struct {
			Key   string
			Value string
		}
		data := testResult{Key: "key", Value: "value"}
		meta := page_meta_dto.PageMetaDto{
			NextPageToken:     new(gofakeit.ID()),
			PreviousPageToken: new(gofakeit.ID()),
			PageSize:          10,
			HasNextPage:       false,
		}
		WriteJsonSuccess(w, data, meta)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response Success[testResult]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		responseMeta := response.Meta
		require.Equal(t, &meta, responseMeta)
		// Need to cast response.Data to map to check content
		require.Equal(t, data.Value, response.Data.Value)
		require.Equal(t, data.Key, response.Data.Key)
	})
	t.Run("Success without meta", func(t *testing.T) {
		w := httptest.NewRecorder()
		type testResult struct {
			Key   string
			Value string
		}
		data := testResult{Key: "key", Value: "value"}
		WriteJsonSuccess(w, data)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response Success[testResult]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		require.Equal(t, data.Key, response.Data.Key)
		require.Equal(t, data.Value, response.Data.Value)
		require.Nil(t, response.Meta)
	})

	t.Run("Success StatusCreated", func(t *testing.T) {
		w := httptest.NewRecorder()
		uuidToReturn := uuid.New()
		WriteJsonCreated(w, uuidToReturn)

		require.Equal(t, http.StatusCreated, w.Code)

		var response Success[uuid.UUID]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		require.Equal(t, uuidToReturn, response.Data)
	})
}
