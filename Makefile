# 📍 Location Demo Project Makefile

.PHONY: help build up down logs ps clean restart \
	backend-build backend-run backend-test backend-tidy \
	frontend-install frontend-dev frontend-build frontend-start

# --- Variables ---
BACKEND_DIR=backend
FRONTEND_DIR=frontend
DOCKER_COMPOSE=docker compose

# --- Help ---
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Docker Targets:"
	@echo "  build            Build all services"
	@echo "  up               Start all services (detached)"
	@echo "  down             Stop and remove all services"
	@echo "  logs             View logs for all services"
	@echo "  ps               List running services"
	@echo "  restart          Restart all services"
	@echo "  clean            Remove all Docker data (containers, volumes, images)"
	@echo ""
	@echo "Backend Targets (Go):"
	@echo "  backend-build    Compile the Go code"
	@echo "  backend-run      Run backend locally (requires DB on 5433)"
	@echo "  backend-test     Run Go tests"
	@echo "  backend-tidy     Run go mod tidy"
	@echo ""
	@echo "Frontend Targets (Next.js):"
	@echo "  frontend-install Install dependencies"
	@echo "  frontend-dev     Start Next.js dev server"
	@echo "  frontend-build   Build Next.js for production"
	@echo "  frontend-start   Start Next.js production server"

# --- Docker ---
build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

ps:
	$(DOCKER_COMPOSE) ps

restart:
	$(DOCKER_COMPOSE) restart

clean:
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

# --- Backend (Local) ---
backend-build:
	cd $(BACKEND_DIR) && go build ./...

backend-run:
	cd $(BACKEND_DIR) && go run ./cmd/api/main.go

backend-test:
	cd $(BACKEND_DIR) && go test ./...

backend-tidy:
	cd $(BACKEND_DIR) && go mod tidy

# --- Frontend (Local) ---
frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev -- -p 3001

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-start:
	cd $(FRONTEND_DIR) && npm run start -- -p 3001
