# Variables
APP_NAME := user-management
MAIN_PATH := ./cmd/main.go
BINARY_PATH := ./bin/$(APP_NAME)
DOCKER_COMPOSE := ./docker-compose.yml

#default
all:lint test build

build:
	@echo "Building Library"
	go build -o $(BINARY_PATH) $(MAIN_PATH)

run:
	@echo "Running Application"
	go run $(MAIN_PATH)

test:
	@echo "Running Tests"
	go test -v ./...

docker-up:
	docker-compose up

docker-down:
	docker-compose down

lint:
	@echo "Running Linter"
	golangci-lint run

