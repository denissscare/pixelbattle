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

DB_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)
include .env
export $(shell sed 's/=.*//' .env)

migrate-up:
	goose -dir ./migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_DSN)" down

migrate-reset:
	goose -dir ./migrations postgres "$(DB_DSN)" reset

migrate-status:
	goose -dir ./migrations postgres "$(DB_DSN)" status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=MigrationName"; \
		exit 1; \
	fi
	goose -dir ./migrations create $(name) sql