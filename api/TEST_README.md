# Go API Test Suite Documentation

This document provides comprehensive information about the test suite for the Grocer-Me Go API.

## Table of Contents

- [Overview](#overview)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Coverage](#test-coverage)
- [Writing Tests](#writing-tests)
- [Integration Testing](#integration-testing)

## Overview

The test suite provides comprehensive coverage of the Go API, including:

- **Handler Tests**: Auth and list management endpoints
- **Middleware Tests**: JWT authentication and CORS
- **Utils Tests**: HTTP utilities, routing, and auth helpers
- **Config Tests**: Configuration initialization

## Test Structure

```
api/
├── handlers/
│   ├── auth_test.go          # Authentication handler tests
│   └── lists_test.go         # List handler tests
├── middleware/
│   ├── jwt_test.go           # JWT middleware tests
│   └── cors_test.go          # CORS middleware tests
├── utils/
│   ├── auth_test.go          # Auth utility tests
│   ├── http_test.go          # HTTP utility tests
│   └── router_test.go        # Router tests
├── config/
│   └── config_test.go        # Configuration tests
└── testutil/
    └── testutil.go           # Test helpers and utilities
```

## Running Tests

### Run All Tests

```bash
make test-api
```

or

```bash
cd api && go test ./... -v
```

### Run Unit Tests Only

These tests don't require a database connection:

```bash
make test-api-unit
```

or

```bash
cd api && go test ./middleware ./utils ./config -v
```

### Run Handler Tests

Handler tests require a test MongoDB instance:

```bash
make test-api-handlers
```

### Run with Coverage Report

```bash
make test-api-coverage
```

This generates:
- `coverage.out`: Coverage data
- `coverage.html`: HTML coverage report (open in browser)

### Run Specific Package

```bash
cd api && go test ./middleware -v
cd api && go test ./utils -v
cd api && go test ./handlers -v
```

### Run Specific Test

```bash
cd api && go test ./middleware -v -run TestJWTAuth_ValidToken
```

### Skip Integration Tests

Many handler tests are marked with `t.Skip()` as they require a database. To run only non-skipped tests:

```bash
cd api && go test ./... -v -short
```

## Test Coverage

### Current Coverage by Package

| Package    | Coverage Type | Description |
|------------|---------------|-------------|
| middleware | Unit Tests    | JWT auth, CORS headers |
| utils      | Unit Tests    | HTTP helpers, router, auth utilities |
| config     | Unit Tests    | Config initialization |
| handlers   | Integration   | API endpoints (requires DB) |

### Generating Coverage Report

```bash
make test-api-coverage
open api/coverage.html
```

## Writing Tests

### Test Utilities

The `testutil` package provides helpful utilities:

```go
// Create authenticated request
req := testutil.CreateAuthenticatedRequest(t, "GET", "/test", nil, userID)

// Create request with JWT token
req := testutil.CreateRequestWithToken(t, "GET", "/test", nil, token)

// Generate test JWT token
token, err := testutil.GenerateTestToken(userID, secret)

// Parse JSON response
var response models.AuthResponse
testutil.ParseJSONResponse(t, w, &response)

// Assert error response
testutil.AssertErrorResponse(t, w, http.StatusUnauthorized)

// Setup test config
testutil.SetupTestConfig(t)
```

### Example Test Structure

```go
func TestMyHandler(t *testing.T) {
    // Setup
    testutil.SetupTestConfig(t)
    
    // Create test data
    userID := "507f1f77bcf86cd799439011"
    
    // Create request
    req := testutil.CreateAuthenticatedRequest(t, "GET", "/test", nil, userID)
    w := httptest.NewRecorder()
    
    // Execute handler
    MyHandler(w, req)
    
    // Assert results
    assert.Equal(t, http.StatusOK, w.Code)
    
    var response MyResponse
    testutil.ParseJSONResponse(t, w, &response)
    assert.Equal(t, expectedValue, response.Field)
}
```

### Test Naming Convention

- Test functions: `Test<FunctionName>_<Scenario>`
- Example: `TestJWTAuth_ValidToken`, `TestHandleSignup_InvalidJSON`

### Table-Driven Tests

For testing multiple scenarios:

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expected    string
        expectError bool
    }{
        {
            name:     "Valid input",
            input:    "test",
            expected: "TEST",
        },
        {
            name:        "Invalid input",
            input:       "",
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := MyFunction(tt.input)
            
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}
```

## Integration Testing

### Setting Up Test Database

For full integration tests, you need a test MongoDB instance:

1. **Using Docker**:

```bash
docker run -d -p 27017:27017 --name mongo-test mongo:latest
```

2. **Set Environment Variables**:

```bash
export MONGODB_URI="mongodb://localhost:27017/grocer-me-test"
export JWT_SECRET="test-secret-for-testing"
```

3. **Run Integration Tests**:

```bash
cd api && go test ./handlers -v
```

### Test Database Best Practices

1. **Isolation**: Each test should clean up after itself
2. **Fixtures**: Use test fixtures for consistent data
3. **Transactions**: Use transactions for atomic test operations
4. **Cleanup**: Always cleanup test data

### Example Integration Test Setup

```go
func TestHandleCreateList_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup test database
    client := setupTestDB(t)
    defer cleanupTestDB(t, client)
    
    // Run test
    // ...
}
```

## Test Environment Variables

Required for integration tests:

```bash
# .env.test
MONGODB_URI=mongodb://localhost:27017/grocer-me-test
JWT_SECRET=test-secret-key-for-testing-purposes-only
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Test API

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mongodb:
        image: mongo:latest
        ports:
          - 27017:27017
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: cd api && go mod download
      
      - name: Run unit tests
        run: make test-api-unit
      
      - name: Run integration tests
        env:
          MONGODB_URI: mongodb://localhost:27017/test
          JWT_SECRET: test-secret
        run: make test-api-handlers
      
      - name: Generate coverage
        run: make test-api-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./api/coverage.out
```

## Troubleshooting

### Common Issues

1. **"JWT_SECRET not set"**
   - Solution: Set `JWT_SECRET` environment variable
   - For tests: `export JWT_SECRET="test-secret"`

2. **"Cannot connect to MongoDB"**
   - Solution: Ensure MongoDB is running
   - For tests: Use `t.Skip()` if database not needed

3. **"Test timeout"**
   - Solution: Increase timeout: `go test ./... -v -timeout 30s`

4. **Import cycle errors**
   - Solution: Check for circular dependencies between packages

## Best Practices

1. **Test Independence**: Each test should be independent
2. **Clear Names**: Use descriptive test names
3. **Setup/Teardown**: Use `t.Cleanup()` for cleanup
4. **Mock External Deps**: Mock external services and databases when possible
5. **Test Edge Cases**: Test error conditions and edge cases
6. **Keep Tests Fast**: Unit tests should be fast; use integration tests sparingly

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

## Contributing

When adding new features:

1. Write tests first (TDD)
2. Ensure existing tests pass
3. Add tests for new functionality
4. Update this documentation if needed
5. Run `make test-api-coverage` to check coverage

## Support

For questions or issues with the test suite, please:

1. Check this documentation
2. Review existing tests for examples
3. Check the troubleshooting section
4. Open an issue on GitHub
