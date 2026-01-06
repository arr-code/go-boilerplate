.PHONY: help migrate-up migrate-down migrate-create sqlc server docker-up docker-down docker-build test clean

# Database connection string for migrations
DB_URL=postgresql://postgres:postgres@localhost:5432/boilerplate_db?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback database migrations"
	@echo "  make migrate-create  - Create new migration (usage: make migrate-create name=migration_name)"
	@echo "  make sqlc            - Generate sqlc code"
	@echo "  make server          - Run the server"
	@echo "  make docker-up       - Start Docker containers"
	@echo "  make docker-down     - Stop Docker containers"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean generated files"

migrate-up:
	@echo "Running migrations..."
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrate-create:
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	@echo "Generating sqlc code..."
	sqlc generate

server:
	@echo "Starting server..."
	go run cmd/api/main.go

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-build:
	@echo "Building Docker image..."
	docker build -t xnoia-go-boilerplate:latest .

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning generated files..."
	rm -rf db/sqlc/*.go
