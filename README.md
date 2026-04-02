# 📍 Location Demo System

A full-stack reference implementation demonstrating how to build a scalable, multi-language location search system using **Clean Architecture**. This project features an intelligent "Water Fall" search strategy that gracefully degrades from local aliases to a free external provider (OpenStreetMap).

## 🌟 Key Features

*   **Intelligent Waterfall Search:** 
    1. Local Alias matching (e.g., query `"sai gon"` matches `Ho Chi Minh City`)
    2. Local Translation matching (e.g., query `"Hồ Chí Minh"` in Vietnamese)
    3. External Fallback (Free OpenStreetMap Nominatim lookup, mapped into DB)
*   **Multi-Language UI:** Seamless server-side rendering for `en`, `vi`, and `ja`.
*   **Clean Architecture (Go):** Strict separation of concerns using decoupled domain models, handlers, services, and repositories.
*   **Hierarchical Locations:** Stores and tracks hierarchical paths (e.g., Country → City → District → Landmark).

## 🛠️ Technology Stack

| Layer | Technologies |
| :--- | :--- |
| **Frontend** | [Next.js 14+](https://nextjs.org/) (App Router), React, Tailwind CSS |
| **Backend** | [Go 1.25+](https://golang.org/), [Gin Web Framework](https://gin-gonic.com/) |
| **Database** | [PostgreSQL 16](https://www.postgresql.org/) |
| **External API**| OpenStreetMap (OSM) Nominatim |
| **Infra** | Docker & Docker Compose |

## 🚀 Quick Start (Recommended)

The easiest way to run the full stack is via Docker. The application is pre-configured to automatically run database migrations and seed data on startup.

**Prerequisites:** 
*   Docker & Docker Compose v2 (`docker compose`)
*   Make

```bash
# 1. Build and start all services (Database, API, Frontend)
make up

# 2. View logs to ensure everything is running smoothly
make logs
```

Once running, you can access the applications locally:
*   **Frontend UI:** [http://localhost:3001](http://localhost:3001)
*   **Backend API:** [http://localhost:8088/health](http://localhost:8088/health)
    *   *Try a search:* `http://localhost:8088/api/v1/locations/search?q=tokyo&lang=en`
*   **PostgreSQL:** `localhost:5433` (User: `postgres` / Pass: `postgres`)

To stop the environment:
```bash
make down
```

## ⌨️ Useful Make Commands

There is a `Makefile` included to help you manage the project efficiently:

| Command | Action |
| :--- | :--- |
| `make build` | Builds all Docker images (API & Frontend) |
| `make up` | Starts all services in the background |
| `make down` | Stops and removes Docker containers and networks |
| `make logs` | Tails the logs for all running Docker services |
| `make clean` | Wipes **ALL** Docker data, including database volumes |
| `make backend-run` | Runs the Go backend locally (without Docker) |
| `make frontend-dev`| Starts Next.js development server locally |

*(Run `make help` to see all available commands).*

## 📖 Documentation

If you want to understand the design decisions, schema, and API structure in-depth, refer to the documentation files in the `docs/` directory:

*   [`docs/gui.md`](./docs/gui.md): High-level system architecture, DB Schema, and project layout.
*   [`docs/how_to_get_gg_map_api_key.md`](./docs/how_to_get_gg_map_api_key.md): Details on external providers. *(Note: The application uses free OpenStreetMap by default, so you don't need Google API keys to get started!)*
*   [`docs/location.md`](./docs/location.md): Foundational database planning/theory notes.
