.PHONY: help install dev build clean test docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing frontend dependencies..."
	cd frontend && npm install
	@echo "Installing backend dependencies..."
	cd backend && go mod tidy
	@echo "Installing AI service dependencies..."
	cd ai-service && pip install -r requirements.txt

dev: ## Start all services in development mode
	@echo "Starting development environment..."
	docker-compose up -d

dev-frontend: ## Start only frontend in development mode
	@echo "Starting frontend..."
	cd frontend && npm run dev

dev-backend: ## Start only backend in development mode
	@echo "Starting backend..."
	cd backend && go run main.go

dev-ai: ## Start only AI service in development mode
	@echo "Starting AI service..."
	cd ai-service && uvicorn main:app --reload --host 0.0.0.0 --port 8000

build: ## Build all services
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Building backend..."
	cd backend && go build -o bin/server main.go
	@echo "Building AI service..."
	cd ai-service && echo "Python service ready"

docker-up: ## Start all services with Docker
	@echo "Starting services with Docker..."
	docker-compose up -d

docker-down: ## Stop all Docker services
	@echo "Stopping Docker services..."
	docker-compose down

docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	docker-compose build

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf frontend/.next
	rm -rf frontend/node_modules
	rm -rf backend/bin
	rm -rf ai-service/__pycache__
	rm -rf ai-service/*.pyc

test: ## Run tests
	@echo "Running frontend tests..."
	cd frontend && npm test
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running AI service tests..."
	cd ai-service && python -m pytest

logs: ## Show logs from all services
	docker-compose logs -f

logs-frontend: ## Show frontend logs
	docker-compose logs -f frontend

logs-backend: ## Show backend logs
	docker-compose logs -f backend

logs-ai: ## Show AI service logs
	docker-compose logs -f ai-service

status: ## Show status of all services
	@echo "Service Status:"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend: http://localhost:8080"
	@echo "AI Service: http://localhost:8000"
	@echo "Swagger Docs: http://localhost:8080/swagger/index.html"
	@echo "AI Service Docs: http://localhost:8000/docs" 