package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bryce-stabenow/grocer-me/testutil"
	"bryce-stabenow/grocer-me/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTAuth_ValidToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create a test user and token
	userID := "507f1f77bcf86cd799439011"
	token, err := testutil.GenerateTestToken(userID, "test-secret-key-for-testing-purposes-only")
	require.NoError(t, err)

	// Create a test handler that checks user ID in context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extractedUserID, ok := utils.GetUserID(r)
		assert.True(t, ok)
		assert.Equal(t, userID, extractedUserID)
		utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
	})

	// Wrap handler with JWT middleware
	handler := JWTAuth(testHandler)

	// Create request with valid token
	req := testutil.CreateRequestWithToken(t, "GET", "/test", nil, token)
	w := httptest.NewRecorder()

	// Execute request
	handler(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_ValidTokenInCookie(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create a test user and token
	userID := "507f1f77bcf86cd799439011"
	token, err := testutil.GenerateTestToken(userID, "test-secret-key-for-testing-purposes-only")
	require.NoError(t, err)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		extractedUserID, ok := utils.GetUserID(r)
		assert.True(t, ok)
		assert.Equal(t, userID, extractedUserID)
		utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
	})

	// Wrap handler with JWT middleware
	handler := JWTAuth(testHandler)

	// Create request with token in cookie
	req := testutil.CreateRequestWithCookie(t, "GET", "/test", nil, "jwt_token", token)
	w := httptest.NewRecorder()

	// Execute request
	handler(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestJWTAuth_NoToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when no token is provided")
	})

	// Wrap handler with JWT middleware
	handler := JWTAuth(testHandler)

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute request
	handler(w, req)

	// Assert response
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	})

	// Wrap handler with JWT middleware
	handler := JWTAuth(testHandler)

	// Create request with invalid token
	req := testutil.CreateRequestWithToken(t, "GET", "/test", nil, "invalid.token.here")
	w := httptest.NewRecorder()

	// Execute request
	handler(w, req)

	// Assert response
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestJWTAuth_TokenWithInvalidSecret(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create token with different secret
	userID := "507f1f77bcf86cd799439011"
	token, err := testutil.GenerateTestToken(userID, "wrong-secret")
	require.NoError(t, err)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with token signed with wrong secret")
	})

	// Wrap handler with JWT middleware
	handler := JWTAuth(testHandler)

	// Create request with token signed with wrong secret
	req := testutil.CreateRequestWithToken(t, "GET", "/test", nil, token)
	w := httptest.NewRecorder()

	// Execute request
	handler(w, req)

	// Assert response
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestExtractUserID_ValidToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create a test user and token
	userID := "507f1f77bcf86cd799439011"
	token, err := testutil.GenerateTestToken(userID, "test-secret-key-for-testing-purposes-only")
	require.NoError(t, err)

	// Create request with valid token
	req := testutil.CreateRequestWithToken(t, "GET", "/test", nil, token)

	// Extract user ID
	extractedUserID, err := ExtractUserID(req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestExtractUserID_NoToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create request without token
	req := httptest.NewRequest("GET", "/test", nil)

	// Extract user ID
	_, err := ExtractUserID(req)

	// Assert
	require.Error(t, err)
}

func TestExtractUserID_InvalidToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create request with invalid token
	req := testutil.CreateRequestWithToken(t, "GET", "/test", nil, "invalid.token")

	// Extract user ID
	_, err := ExtractUserID(req)

	// Assert
	require.Error(t, err)
}

func TestExtractUserID_TokenInCookie(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Create a test user and token
	userID := "507f1f77bcf86cd799439011"
	token, err := testutil.GenerateTestToken(userID, "test-secret-key-for-testing-purposes-only")
	require.NoError(t, err)

	// Create request with token in cookie
	req := testutil.CreateRequestWithCookie(t, "GET", "/test", nil, "jwt_token", token)

	// Extract user ID
	extractedUserID, err := ExtractUserID(req)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}
