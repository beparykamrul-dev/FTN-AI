.PHONY: help build run test clean dev production logs stop

# Variables
IMAGE_NAME=ftn-ai
CONTAINER_NAME=ftn-backend
BRANCH=$$(git rev-parse --abbrev-ref HEAD)
COMMIT=$$(git rev-parse --short HEAD)

help:
	@echo "FTN-AI Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev           - Start development environment"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make docker-stop   - Stop Docker container"
	@echo ""
	@echo "Production:"
	@echo "  make production    - Production setup"
	@echo "  make deploy        - Deploy to production"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make logs          - View application logs"
	@echo "  make stop          - Stop all services"

dev:
	@echo "Starting development environment..."
	docker-compose up -d
	@echo "Services running. Access API at http://localhost:8080"

build:
	@echo "Building application..."
	cd backend && go build -o ftn-backend main.go
	@echo "Build complete"

run: build
	@echo "Running application..."
	./backend/ftn-backend

test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):latest .
	@echo "Docker image built: $(IMAGE_NAME):latest"

docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(CONTAINER_NAME) -p 8080:8080 $(IMAGE_NAME):latest

docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

production:
	@echo "Setting up production environment..."
	docker-compose -f docker-compose.yml up -d
	@echo "Production environment ready"

clean:
	@echo "Cleaning up..."
	rm -rf backend/ftn-backend
	rm -rf coverage.out
	@echo "Clean complete"

logs:
	docker-compose logs -f

stop:
	@echo "Stopping all services..."
	docker-compose down
	@echo "All services stopped"