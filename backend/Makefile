GO := go
LINTER := golangci-lint
DOCKER_COMPOSE := docker-compose
IMAGE_NAME := pixel-battle

all: lint build

run:
	@echo "Running project"
	$(DOCKER_COMPOSE) up -d

lint:
	@echo "Running linters..."
	$(LINTER) run

build:
	@echo "Building Docker image with docker-compose..."
	$(DOCKER_COMPOSE) up -d --build

stop:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

