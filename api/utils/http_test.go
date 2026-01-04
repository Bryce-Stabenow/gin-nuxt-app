package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
	}{
		{
			name:       "Simple map response",
			statusCode: http.StatusOK,
			data:       map[string]string{"message": "success"},
		},
		{
			name:       "Array response",
			statusCode: http.StatusOK,
			data:       []string{"item1", "item2", "item3"},
		},
		{
			name:       "Struct response",
			statusCode: http.StatusCreated,
			data: struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{
				ID:   "123",
				Name: "Test",
			},
		},
		{
			name:       "Nil data",
			statusCode: http.StatusNoContent,
			data:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			JSONResponse(w, tt.statusCode, tt.data)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.data != nil {
				var result interface{}
				err := json.NewDecoder(w.Body).Decode(&result)
				require.NoError(t, err)
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
	}{
		{
			name:       "Bad request error",
			statusCode: http.StatusBadRequest,
			message:    "Invalid request body",
		},
		{
			name:       "Unauthorized error",
			statusCode: http.StatusUnauthorized,
			message:    "Authentication required",
		},
		{
			name:       "Not found error",
			statusCode: http.StatusNotFound,
			message:    "Resource not found",
		},
		{
			name:       "Internal server error",
			statusCode: http.StatusInternalServerError,
			message:    "Something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ErrorResponse(w, tt.statusCode, tt.message)

			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response map[string]string
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, tt.message, response["error"])
		})
	}
}

func TestDecodeJSON(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	t.Run("Valid JSON", func(t *testing.T) {
		jsonData := `{"name":"John Doe","email":"john@example.com","age":30}`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonData))

		var result TestStruct
		err := DecodeJSON(req, &result)

		require.NoError(t, err)
		assert.Equal(t, "John Doe", result.Name)
		assert.Equal(t, "john@example.com", result.Email)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		jsonData := `{"name":"John Doe","email":}`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(jsonData))

		var result TestStruct
		err := DecodeJSON(req, &result)

		require.Error(t, err)
	})

	t.Run("Empty body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))

		var result TestStruct
		err := DecodeJSON(req, &result)

		require.Error(t, err)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("User ID exists in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		userID := "507f1f77bcf86cd799439011"
		req = SetUserID(req, userID)

		extractedUserID, ok := GetUserID(req)

		assert.True(t, ok)
		assert.Equal(t, userID, extractedUserID)
	})

	t.Run("User ID not in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)

		_, ok := GetUserID(req)

		assert.False(t, ok)
	})
}

func TestSetUserID(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	userID := "507f1f77bcf86cd799439011"

	newReq := SetUserID(req, userID)

	extractedUserID, ok := GetUserID(newReq)
	assert.True(t, ok)
	assert.Equal(t, userID, extractedUserID)

	// Verify original request is not modified
	_, ok = GetUserID(req)
	assert.False(t, ok)
}

func TestGetPathParam(t *testing.T) {
	t.Run("Path param exists", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/lists/123", nil)
		params := map[string]string{"id": "123", "name": "test"}
		req = SetPathParams(req, params)

		id := GetPathParam(req, "id")
		name := GetPathParam(req, "name")

		assert.Equal(t, "123", id)
		assert.Equal(t, "test", name)
	})

	t.Run("Path param not exists", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/lists/123", nil)
		params := map[string]string{"id": "123"}
		req = SetPathParams(req, params)

		nonExistent := GetPathParam(req, "nonexistent")

		assert.Equal(t, "", nonExistent)
	})

	t.Run("No path params in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)

		value := GetPathParam(req, "id")

		assert.Equal(t, "", value)
	})
}

func TestSetPathParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	params := map[string]string{
		"id":   "123",
		"name": "test",
	}

	newReq := SetPathParams(req, params)

	id := GetPathParam(newReq, "id")
	name := GetPathParam(newReq, "name")
	assert.Equal(t, "123", id)
	assert.Equal(t, "test", name)

	// Verify original request is not modified
	origID := GetPathParam(req, "id")
	assert.Equal(t, "", origID)
}

func TestSetCookie(t *testing.T) {
	tests := []struct {
		name     string
		name_    string
		value    string
		maxAge   int
		path     string
		domain   string
		secure   bool
		httpOnly bool
	}{
		{
			name:     "Basic cookie",
			name_:    "session",
			value:    "abc123",
			maxAge:   3600,
			path:     "/",
			domain:   "",
			secure:   false,
			httpOnly: true,
		},
		{
			name:     "Secure cookie",
			name_:    "auth_token",
			value:    "xyz789",
			maxAge:   86400,
			path:     "/api",
			domain:   "example.com",
			secure:   true,
			httpOnly: true,
		},
		{
			name:     "Non-HTTP-only cookie",
			name_:    "preference",
			value:    "dark_mode",
			maxAge:   604800,
			path:     "/",
			domain:   "",
			secure:   false,
			httpOnly: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			SetCookie(w, tt.name_, tt.value, tt.maxAge, tt.path, tt.domain, tt.secure, tt.httpOnly)

			cookies := w.Result().Cookies()
			require.Len(t, cookies, 1)

			cookie := cookies[0]
			assert.Equal(t, tt.name_, cookie.Name)
			assert.Equal(t, tt.value, cookie.Value)
			assert.Equal(t, tt.maxAge, cookie.MaxAge)
			assert.Equal(t, tt.path, cookie.Path)
			assert.Equal(t, tt.domain, cookie.Domain)
			assert.Equal(t, tt.secure, cookie.Secure)
			assert.Equal(t, tt.httpOnly, cookie.HttpOnly)
		})
	}
}
