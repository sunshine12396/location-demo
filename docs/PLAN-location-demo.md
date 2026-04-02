# Implementation Plan: Location Demo System

## Overview
Build a full-stack location demo with Go backend (Clean Architecture) and Next.js frontend.

## Phases

### Phase 1: Backend Foundation ✅
- [ ] Go module init + project structure
- [ ] Domain models & interfaces (`/internal/domain`)
- [ ] Database migrations (`/migrations`)
- [ ] Seed data (`/migrations/seed.sql`)
- [ ] PostgreSQL repository implementation
- [ ] Location service (business logic with waterfall search)
- [ ] HTTP handlers (Gin router)
- [ ] Config loading (env-based)
- [ ] `main.go` wiring

### Phase 2: Frontend Foundation
- [ ] Next.js App Router init
- [ ] Search page (client component with debounce)
- [ ] Location detail page (server component)
- [ ] Shared components (SearchBar, LocationCard)
- [ ] API client layer

### Phase 3: Polish
- [ ] Docker Compose (PostgreSQL + Backend + Frontend)
- [ ] README with setup instructions

## File Map

```
backend/
  go.mod
  cmd/api/main.go
  internal/
    config/config.go
    domain/location.go
    location/
      handler.go
      service.go
      repository.go
  migrations/
    001_create_tables.up.sql
    001_create_tables.down.sql
    seed.sql

frontend/
  package.json
  next.config.js
  app/
    layout.tsx
    page.tsx              # Search page
    location/[id]/page.tsx
  components/
    SearchBar.tsx
    LocationCard.tsx
  lib/
    api.ts

docker-compose.yml
```
