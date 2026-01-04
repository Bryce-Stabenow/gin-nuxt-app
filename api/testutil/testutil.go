package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bryce-stabenow/grocer-me/config"
	"bryce-stabenow/grocer-me/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateTestUser creates a test user and returns the user
func CreateTestUser(t *testing.T) *models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), 10)
	require.NoError(t, err)

	now := time.Now()
	return &models.User{
		ID:           primitive.NewObjectID(),
		Email:        "test@example.com",
		Username:     "test@example.com",
		PasswordHash: string(hashedPassword),
		Profile: &models.Profile{
			FirstName: "Test",
			LastName:  "User",
			AvatarURL: "https://example.com/avatar.png",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateTestList creates a test list
func CreateTestList(t *testing.T, userID primitive.ObjectID) *models.List {
	now := time.Now()
	return &models.List{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        "Test List",
		Description: "A test list",
		Items: []models.ListItem{
			{
				Name:     "Test Item",
				Quantity: 1,
				Checked:  false,
				Details:  "Test details",
				AddedBy:  userID,
				AddedAt:  now,
			},
		},
		SharedWith: []primitive.ObjectID{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// GenerateTestToken generates a JWT token for testing
func GenerateTestToken(userID string, secret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// CreateAuthenticatedRequest creates an HTTP request with user ID in context
// Note: This requires the utils package to be imported by the calling test
func CreateAuthenticatedRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	return req
}

// CreateRequestWithToken creates an HTTP request with JWT token in header
func CreateRequestWithToken(t *testing.T, method, url string, body interface{}, token string) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req
}

// ParseJSONResponse parses JSON response from ResponseRecorder
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	err := json.NewDecoder(w.Body).Decode(v)
	require.NoError(t, err)
}

// SetupTestConfig sets up test configuration
func SetupTestConfig(t *testing.T) {
	config.JWTSecret = "test-secret-key-for-testing-purposes-only"
}

// CreateRequestWithPathParams creates a request (path params should be set by caller)
func CreateRequestWithPathParams(t *testing.T, method, url string, body interface{}) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	return req
}

// SetupMockDB sets up a mock database for testing
func SetupMockDB(t *testing.T) {
	// Note: In a real scenario, you would use a test MongoDB instance
	// For now, this is a placeholder that would require proper mocking
	// or integration with a test database
}

// CleanupMockDB cleans up mock database
func CleanupMockDB(t *testing.T) {
	// Cleanup logic would go here
}

// AssertErrorResponse checks if response contains an error message
func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	require.Equal(t, expectedStatus, w.Code)
	
	var response map[string]interface{}
	ParseJSONResponse(t, w, &response)
	
	_, hasError := response["error"]
	require.True(t, hasError, "Response should contain error field")
}

// AssertSuccessResponse checks if response is successful
func AssertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) {
	require.Equal(t, expectedStatus, w.Code)
}

// CreateRequestWithCookie creates a request with a cookie
func CreateRequestWithCookie(t *testing.T, method, url string, body interface{}, cookieName, cookieValue string) *http.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	if cookieName != "" && cookieValue != "" {
		req.AddCookie(&http.Cookie{
			Name:  cookieName,
			Value: cookieValue,
		})
	}

	return req
}

// GetMockContext returns a mock context with timeout
func GetMockContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
