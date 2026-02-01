.PHONY: build run stop restart logs clean test help dev swagger

APP_NAME=kasir-api
DOCKER_IMAGE=$(APP_NAME):latest

build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

dev:
	@echo "Running locally..."
	go run ./cmd/api

swagger:
	@echo "Generating swagger docs..."
	swag init -g cmd/api/main.go -o docs

run:
	@echo "Starting application..."
	docker-compose up -d

stop:
	@echo "Stopping application..."
	docker-compose down

restart: stop run

logs:
	@echo "Showing logs (Ctrl+C to exit)..."
	docker-compose logs -f

clean:
	@echo "Cleaning up..."
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true

test:
	@echo "Testing health endpoint..."
	@curl -f http://localhost:8080/health && echo "\n✓ Health check passed!" || echo "\n✗ Health check failed!"

test-all:
	@echo "Testing all endpoints..."
	@echo "\n1. Health check:"
	@curl -s http://localhost:8080/health | python3 -m json.tool || echo "Failed"
	@echo "\n2. Get all products:"
	@curl -s http://localhost:8080/api/product | python3 -m json.tool || echo "Failed"
	@echo "\n3. Get product by ID (1):"
	@curl -s http://localhost:8080/api/product/1 | python3 -m json.tool || echo "Failed"

help:
	@echo "Available commands:"
	@echo "  make build     - Build Docker image"
	@echo "  make dev       - Run application locally"
	@echo "  make swagger   - Generate swagger docs"
	@echo "  make run       - Start application with docker-compose"
	@echo "  make stop      - Stop application"
	@echo "  make restart   - Restart application"
	@echo "  make logs      - View application logs"
	@echo "  make clean     - Remove containers and images"
	@echo "  make test      - Test health endpoint"
	@echo "  make test-all  - Test all API endpoints"
	@echo "  make help      - Show this help message"
