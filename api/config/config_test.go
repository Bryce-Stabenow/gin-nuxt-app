package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_WithoutJWTSecret(t *testing.T) {
	// Clear JWT_SECRET
	os.Unsetenv("JWT_SECRET")

	// This would normally call log.Fatal, which we can't test directly
	// In a production test suite, you might use a testing framework that
	// can capture os.Exit calls, or refactor Init to return an error
	// For now, we'll skip this test
	t.Skip("Skipping test that would call log.Fatal")
}

func TestInit_WithJWTSecret(t *testing.T) {
	// Set JWT_SECRET
	testSecret := "test-jwt-secret-12345"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Call Init
	Init()

	// Assert JWT secret is set
	assert.Equal(t, testSecret, JWTSecret)
}

func TestSetMongoClient(t *testing.T) {
	// Note: SetMongoClient requires a valid client to call Database()
	// This test would need a real or mock MongoDB connection
	t.Skip("Requires valid MongoDB client for testing")
}

func TestGetMongoURI_WithoutURI(t *testing.T) {
	// Clear MONGODB_URI
	os.Unsetenv("MONGODB_URI")

	// This would normally call log.Fatal, which we can't test directly
	t.Skip("Skipping test that would call log.Fatal")
}

func TestGetMongoURI_WithURI(t *testing.T) {
	// Set MONGODB_URI
	testURI := "mongodb://localhost:27017/test"
	os.Setenv("MONGODB_URI", testURI)
	defer os.Unsetenv("MONGODB_URI")

	// Call GetMongoURI
	uri := GetMongoURI()

	// Assert URI is returned
	assert.Equal(t, testURI, uri)
}

func TestGlobalVariables(t *testing.T) {
	// Test that global variables can be set and retrieved
	testSecret := "my-test-secret"
	JWTSecret = testSecret

	assert.Equal(t, testSecret, JWTSecret)
}

func TestInit_LoadsEnvironmentVariables(t *testing.T) {
	// Set environment variable
	testSecret := "env-test-secret"
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	// Clear current value
	JWTSecret = ""

	// Call Init
	Init()

	// Assert value was loaded from environment
	assert.Equal(t, testSecret, JWTSecret)
}

func TestGetMongoURI_ReturnsCorrectValue(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "Standard MongoDB URI",
			envValue: "mongodb://localhost:27017",
			expected: "mongodb://localhost:27017",
		},
		{
			name:     "MongoDB Atlas URI",
			envValue: "mongodb+srv://user:pass@cluster.mongodb.net/dbname",
			expected: "mongodb+srv://user:pass@cluster.mongodb.net/dbname",
		},
		{
			name:     "MongoDB with authentication",
			envValue: "mongodb://admin:password@localhost:27017/admin",
			expected: "mongodb://admin:password@localhost:27017/admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("MONGODB_URI", tt.envValue)
			defer os.Unsetenv("MONGODB_URI")

			uri := GetMongoURI()

			assert.Equal(t, tt.expected, uri)
		})
	}
}
