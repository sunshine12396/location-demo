# 📍 Location Demo Project Makefile

.PHONY: help build up down logs ps clean restart \
	install \
	be-build be-dev be-test be-tidy \
	fe-install fe-dev fe-build fe-start \
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
	@echo "Project Setup:"
	@echo "  install          Install all dependencies for FE & BE"
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
	@echo "Backend Targets (Local):"
	@echo "  be-dev           Run backend locally (requires DB on 5432)"
	@echo "  be-test          Run Go tests"
	@echo "  be-tidy          Run go mod tidy"
	@echo ""
	@echo "Frontend Targets (Local):"
	@echo "  fe-dev           Start Next.js dev server"
	@echo "  fe-build         Build Next.js for production"
	@echo ""

# --- Project Setup ---
install: fe-install be-tidy

# --- Docker ---
build:
	$(DOCKER_COMPOSE) build

up:
	$(DOCKER_COMPOSE) up -d

db-up:
	$(DOCKER_COMPOSE) up -d db

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
be-build:
	cd $(BACKEND_DIR) && go build ./...

be-dev:
	cd $(BACKEND_DIR) && go run ./cmd/api/main.go

be-test:
	cd $(BACKEND_DIR) && go test ./...

be-tidy:
	cd $(BACKEND_DIR) && go mod tidy

# Aliases
backend-build: be-build
backend-run: be-dev
backend-test: be-test
backend-tidy: be-tidy

# --- Frontend (Local) ---
fe-install:
	cd $(FRONTEND_DIR) && npm install

fe-dev:
	cd $(FRONTEND_DIR) && npm run dev -- -p 3001

fe-build:
	cd $(FRONTEND_DIR) && npm run build

fe-start:
	cd $(FRONTEND_DIR) && npm run start -- -p 3001

# Aliases
frontend-install: fe-install
frontend-dev: fe-dev
frontend-build: fe-build
frontend-start: fe-start
