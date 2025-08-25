# Makefile for Go service with Redis and PostgreSQL

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the Go application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start all services with Docker Compose"
	@echo "  docker-down  - Stop all services"
	@echo "  docker-logs  - View logs from all services"
	@echo "  api-test     - Test API endpoints"

# Build the Go application
build:
	go mod download
	go build -o bin/main .

# Run the application locally (requires local Redis and PostgreSQL)
run:
	go run .

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t user-service .

# Start all services with Docker Compose
docker-up:
	docker-compose up -d

# Stop all services
docker-down:
	docker-compose down

# View logs from all services
docker-logs:
	docker-compose logs -f

# Test API endpoints (requires the service to be running)
api-test:
	@echo "Testing API endpoints..."
	@echo "1. Health check:"
	curl -s http://localhost:8080/health | jq .
	@echo "\n2. Create user:"
	curl -s -X POST http://localhost:8080/api/v1/users \
		-H "Content-Type: application/json" \
		-d '{"name":"Test User","email":"test@example.com"}' | jq .
	@echo "\n3. Get all users:"
	curl -s http://localhost:8080/api/v1/users | jq .
	@echo "\n4. Get user by ID (should be cached):"
	curl -s http://localhost:8080/api/v1/users/1 -I | grep X-Cache || echo "No cache header"
	curl -s http://localhost:8080/api/v1/users/1 | jq .
