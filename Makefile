.PHONY: build up down logs clean test help

.SILENT:

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

clean:
	docker-compose down -v

lint:
	golangci-lint run

run: build up
	@echo "Service is starting..."

restart: down run

health:
	curl -f http://localhost:8080/health || echo "❌ Service is not ready"

help:
	@echo "Quick start: make run"
	@echo "Stop: make down"
	@echo "View logs: make logs"
	@echo "Run tests: make test"

.DEFAULT_GOAL := run