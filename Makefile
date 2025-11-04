.PHONY: help install run build test clean docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make install    - Install Go dependencies"
	@echo "  make run        - Run the backend server"
	@echo "  make build      - Build the backend binary"
	@echo "  make test       - Run tests"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make docker-up  - Start Docker containers"
	@echo "  make docker-down- Stop Docker containers"

install:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy

run:
	@echo "Starting backend server..."
	go run main.go

build:
	@echo "Building backend..."
	go build -o bin/xivercrm main.go

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf data/*.db

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

frontend-dev:
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

frontend-build:
	@echo "Building frontend..."
	cd frontend && npm run build

