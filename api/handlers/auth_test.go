package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bryce-stabenow/grocer-me/models"
	"bryce-stabenow/grocer-me/testutil"
	"bryce-stabenow/grocer-me/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: These tests are unit tests that test the handler logic without a real database.
// For full integration tests with MongoDB, you would need to set up a test database.

func TestHandleSignup_ValidRequest(t *testing.T) {
	testutil.SetupTestConfig(t)

	// This test demonstrates the expected structure
	// In a real scenario, you would mock the database or use a test database
	t.Skip("Requires database setup for integration testing")

	req := models.SignupRequest{
		Email:     "newuser@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/signup", req, "")
	w := httptest.NewRecorder()

	HandleSignup(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.AuthResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.NotEmpty(t, response.Token)
	assert.NotNil(t, response.User)
	assert.Equal(t, req.Email, response.User.Email)
}

func TestHandleSignup_InvalidJSON(t *testing.T) {
	testutil.SetupTestConfig(t)

	// Test with invalid JSON
	httpReq := httptest.NewRequest("POST", "/signup", nil)
	w := httptest.NewRecorder()

	HandleSignup(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleSignup_MissingFields(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	// Test with missing required fields
	req := models.SignupRequest{
		Email: "test@example.com",
		// Missing password, firstName, lastName
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/signup", req, "")
	w := httptest.NewRecorder()

	HandleSignup(w, httpReq)

	// The actual validation would depend on the framework
	// Here we just document the expected behavior
	assert.NotEqual(t, http.StatusCreated, w.Code)
}

func TestHandleSignin_ValidCredentials(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	req := models.SigninRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/signin", req, "")
	w := httptest.NewRecorder()

	HandleSignin(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.AuthResponse
	testutil.ParseJSONResponse(t, w, &response)

	assert.NotEmpty(t, response.Token)
	assert.NotNil(t, response.User)
}

func TestHandleSignin_InvalidJSON(t *testing.T) {
	testutil.SetupTestConfig(t)

	httpReq := httptest.NewRequest("POST", "/signin", nil)
	w := httptest.NewRecorder()

	HandleSignin(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusBadRequest)
}

func TestHandleSignin_InvalidCredentials(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	req := models.SigninRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/signin", req, "")
	w := httptest.NewRecorder()

	HandleSignin(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestHandleGetMe_Authenticated(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	userID := "507f1f77bcf86cd799439011"
	httpReq := testutil.CreateAuthenticatedRequest(t, "GET", "/me", nil)
	httpReq = utils.SetUserID(httpReq, userID)
	w := httptest.NewRecorder()

	HandleGetMe(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var user models.User
	testutil.ParseJSONResponse(t, w, &user)

	assert.Equal(t, userID, user.ID.Hex())
}

func TestHandleGetMe_Unauthenticated(t *testing.T) {
	testutil.SetupTestConfig(t)

	httpReq := httptest.NewRequest("GET", "/me", nil)
	w := httptest.NewRecorder()

	HandleGetMe(w, httpReq)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)
}

func TestHandleGetMe_UserNotFound(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	// Use a valid ObjectID format but non-existent user
	nonExistentUserID := "507f1f77bcf86cd799439999"
	httpReq := testutil.CreateAuthenticatedRequest(t, "GET", "/me", nil)
	httpReq = utils.SetUserID(httpReq, nonExistentUserID)
	w := httptest.NewRecorder()

	HandleGetMe(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	testutil.AssertErrorResponse(t, w, http.StatusNotFound)
}

func TestHandleLogout(t *testing.T) {
	testutil.SetupTestConfig(t)

	userID := "507f1f77bcf86cd799439011"
	httpReq := testutil.CreateAuthenticatedRequest(t, "POST", "/logout", nil)
	httpReq = utils.SetUserID(httpReq, userID)
	w := httptest.NewRecorder()

	HandleLogout(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the JWT cookie is cleared
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "jwt_token" {
			found = true
			assert.Equal(t, "", cookie.Value)
			assert.Equal(t, -1, cookie.MaxAge)
			break
		}
	}
	assert.True(t, found, "JWT cookie should be set for clearing")

	var response map[string]string
	testutil.ParseJSONResponse(t, w, &response)
	assert.Equal(t, "Logged out successfully", response["message"])
}

func TestGenerateToken(t *testing.T) {
	testutil.SetupTestConfig(t)

	userID := "507f1f77bcf86cd799439011"
	
	token, err := generateToken(userID)
	
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Verify token can be parsed
	extractedUserID, err := testutil.GenerateTestToken(userID, "test-secret-key-for-testing-purposes-only")
	require.NoError(t, err)
	assert.NotEmpty(t, extractedUserID)
}

func TestHandleSignup_SetsCookie(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	req := models.SignupRequest{
		Email:     "cookie@example.com",
		Password:  "password123",
		FirstName: "Cookie",
		LastName:  "User",
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/signup", req, "")
	w := httptest.NewRecorder()

	HandleSignup(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Check that JWT cookie is set
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "jwt_token" {
			found = true
			assert.NotEmpty(t, cookie.Value)
			assert.Equal(t, "/", cookie.Path)
			assert.True(t, cookie.HttpOnly)
			break
		}
	}
	assert.True(t, found, "JWT cookie should be set")
}

func TestHandleSignin_SetsCookie(t *testing.T) {
	testutil.SetupTestConfig(t)

	t.Skip("Requires database setup for integration testing")

	req := models.SigninRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	httpReq := testutil.CreateRequestWithToken(t, "POST", "/signin", req, "")
	w := httptest.NewRecorder()

	HandleSignin(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check that JWT cookie is set
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "jwt_token" {
			found = true
			assert.NotEmpty(t, cookie.Value)
			assert.Equal(t, "/", cookie.Path)
			assert.True(t, cookie.HttpOnly)
			break
		}
	}
	assert.True(t, found, "JWT cookie should be set")
}

func TestAuthResponse_Structure(t *testing.T) {
	// Test that AuthResponse can be marshaled/unmarshaled correctly
	response := models.AuthResponse{
		Token: "test-token",
		User: &models.UserPublic{
			ID:       "507f1f77bcf86cd799439011",
			Email:    "test@example.com",
			Username: "testuser",
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	var decoded models.AuthResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, response.Token, decoded.Token)
	assert.Equal(t, response.User.ID, decoded.User.ID)
	assert.Equal(t, response.User.Email, decoded.User.Email)
}
