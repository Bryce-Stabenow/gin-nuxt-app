run-api:
	cd api && go run .

run-web:
	cd web && npm run dev

# Test commands
test-api:
	cd api && go test ./... -v

test-api-coverage:
	cd api && go test ./... -v -coverprofile=coverage.out
	cd api && go tool cover -html=coverage.out -o coverage.html

test-api-short:
	cd api && go test ./... -v -short

test-api-unit:
	cd api && go test ./middleware ./utils ./config -v

test-api-handlers:
	cd api && go test ./handlers -v

# Run all tests (requires test database)
test-all:
	$(MAKE) test-api

# Clean test artifacts
clean-test:
	cd api && rm -f coverage.out coverage.html