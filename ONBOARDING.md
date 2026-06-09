# Onboarding Guide

A step-by-step guide to understanding, configuring, and taking ownership of the **bkit-api-template** for your own project.

---

## Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [Repository Layout](#2-repository-layout)
3. [Configuration — What You Can Change](#3-configuration--what-you-can-change)
   - [Option A: Edit config.yaml](#option-a-edit-configyaml)
   - [Option B: Edit environment variables in compose.yml](#option-b-edit-environment-variables-in-composeyml)
   - [Full Configuration Reference](#full-configuration-reference)
4. [Code — What You Can Change](#4-code--what-you-can-change)
5. [Running the Project](#5-running-the-project)
6. [Taking Over for Your Real Project](#6-taking-over-for-your-real-project)

---

## 1. Prerequisites

| Tool | Minimum Version | Purpose |
|---|---|---|
| **Go** | 1.25+ | Language runtime |
| **Docker & Docker Compose** | Latest | Containerized dev environment |
| **Make** | Any | Task automation |
| **Node.js / npm** | 18+ | Tailwind CSS CLI (installed automatically inside Docker) |

> **Tip:** If you run the full stack via `docker compose up`, all Go tools (`templ`, `air`, `templui`) and Node dependencies are installed inside the container automatically — you only need Docker.

---

## 2. Repository Layout

```text
bkit-api-template/
├── cmd/
│   └── main.go                  # ← Application entrypoint
├── internal/
│   ├── config/
│   │   └── config.go            # ← Config loader (YAML + env var merging)
│   ├── server/
│   │   └── server.go            # ← HTTP server (Echo v4, routes, middleware)
│   └── worker/
│       └── worker.go            # ← Background worker loop
├── views/
│   ├── pages/
│   │   ├── home.templ           # ← Dashboard page template
│   │   └── layout.templ         # ← Base HTML layout
│   ├── components/              # ← Auto-generated templUI components (gitignored)
│   └── utils/                   # ← Auto-generated templUI utilities (gitignored)
├── assets/
│   ├── css/input.css            # ← Tailwind v4 entry point
│   └── js/lib/htmx.js           # ← HTMX library
├── config.yaml                  # ← Default configuration values
├── compose.yml                  # ← Docker Compose (env var overrides live here)
├── Dockerfile                   # ← Production multi-stage build
├── Dockerfile.dev               # ← Development container with hot-reload
├── Makefile                     # ← Build & dev task automation
└── .air.toml                    # ← Hot-reload watcher config
```

---

## 3. Configuration — What You Can Change

There are **two ways** to configure the application. They are merged at startup — environment variables always take precedence over `config.yaml`.

### Option A: Edit `config.yaml`

This is the default configuration file. Values here are the baseline for every environment.

```yaml
# config.yaml

server:
  port: 8080                     # Port the HTTP server listens on

telemetry:
  enabled: true                  # Master toggle for all telemetry
  service_name: "bkit-api-service"  # Identifies your service in traces/metrics
  environment: "development"     # development | staging | production
  enable_trace: false            # Enable distributed tracing (OTLP)
  enable_metrics: true           # Enable Prometheus metrics
  metric_mode: "pull"            # "pull" = expose /metrics endpoint
  enable_stdout: false           # Print telemetry to stdout (debugging)

kv:
  enabled: true                  # Enable key-value cache
  driver: "local"                # Options: local, redis, sqlite
  redis:
    addr: "redis:6379"           # Redis host:port
    db: 0                        # Redis database index
    prefix: "bkit"               # Key prefix in Redis

db:
  enabled: true                  # Enable database
  driver: "postgres"             # Options: postgres, mariadb, sqlite
  postgres:
    host: "postgres"
    port: 5432
    user: "postgres"
    password: "postgrespassword"
    dbname: "bkit_db"
    sslmode: "disable"
  mariadb:
    host: "mariadb"
    port: 3306
    user: "bkit_user"
    password: "mariadbpassword"
    dbname: "bkit_db"
  sqlite:
    path: "data.db"
  max_open_conns: 10             # Max open DB connections
  max_idle_conns: 5              # Max idle DB connections
```

### Option B: Edit environment variables in `compose.yml`

Environment variables in `compose.yml` override matching keys from `config.yaml`. They use `SCREAMING_SNAKE_CASE` and map to the nested YAML structure by replacing dots/nesting with underscores.

**Mapping rule:** `section.subsection.key` → `SECTION_SUBSECTION_KEY`

```yaml
# compose.yml → services → app → environment

environment:
  # Server
  - SERVER_PORT=8080

  # Telemetry
  - TELEMETRY_ENABLED=true
  - TELEMETRY_SERVICE_NAME=bkit-api-service
  - TELEMETRY_ENVIRONMENT=development
  - TELEMETRY_ENABLE_TRACE=false
  - TELEMETRY_ENABLE_METRICS=true
  - TELEMETRY_METRIC_MODE=pull
  - TELEMETRY_ENABLE_STDOUT=false

  # Key-Value Store
  - KV_ENABLED=true
  - KV_DRIVER=redis                  # local | redis | sqlite
  - KV_REDIS_ADDR=redis:6379
  - KV_REDIS_DB=0
  - KV_REDIS_PREFIX=bkit

  # Database
  - DB_ENABLED=true
  - DB_DRIVER=postgres               # postgres | mariadb | sqlite
  - DB_POSTGRES_HOST=postgres
  - DB_POSTGRES_PORT=5432
  - DB_POSTGRES_USER=postgres
  - DB_POSTGRES_PASSWORD=postgrespassword
  - DB_POSTGRES_DBNAME=bkit_db
  - DB_POSTGRES_SSLMODE=disable
  - DB_MARIADB_HOST=mariadb
  - DB_MARIADB_PORT=3306
  - DB_MARIADB_USER=bkit_user
  - DB_MARIADB_PASSWORD=mariadbpassword
  - DB_MARIADB_DBNAME=bkit_db
  - DB_SQLITE_PATH=data.db
  - DB_MAX_OPEN_CONNS=10
  - DB_MAX_IDLE_CONNS=5
```

> **When to use which?**
> - Use `config.yaml` for sensible defaults that you commit to the repo.
> - Use `compose.yml` environment variables for values that change per deployment or contain secrets you don't want in version control.
> - Both approaches work together — env vars override matching YAML keys at startup.

### Full Configuration Reference

| Config Key (dot-path) | Env Variable | Type | Default | Description |
|---|---|---|---|---|
| `server.port` | `SERVER_PORT` | int | `8080` | HTTP server listen port |
| `telemetry.enabled` | `TELEMETRY_ENABLED` | bool | `true` | Master telemetry toggle |
| `telemetry.service_name` | `TELEMETRY_SERVICE_NAME` | string | `bkit-api-service` | Service name in traces/metrics |
| `telemetry.environment` | `TELEMETRY_ENVIRONMENT` | string | `development` | Environment label |
| `telemetry.enable_trace` | `TELEMETRY_ENABLE_TRACE` | bool | `false` | Enable OTLP tracing |
| `telemetry.enable_metrics` | `TELEMETRY_ENABLE_METRICS` | bool | `true` | Enable Prometheus metrics |
| `telemetry.metric_mode` | `TELEMETRY_METRIC_MODE` | string | `pull` | Metric export mode |
| `telemetry.enable_stdout` | `TELEMETRY_ENABLE_STDOUT` | bool | `false` | Print telemetry to stdout |
| `kv.enabled` | `KV_ENABLED` | bool | `true` | Enable key-value store |
| `kv.driver` | `KV_DRIVER` | string | `local` | KV driver: `local`, `redis`, `sqlite` |
| `kv.redis.addr` | `KV_REDIS_ADDR` | string | `redis:6379` | Redis address |
| `kv.redis.db` | `KV_REDIS_DB` | int | `0` | Redis database index |
| `kv.redis.prefix` | `KV_REDIS_PREFIX` | string | `bkit` | Redis key prefix |
| `db.enabled` | `DB_ENABLED` | bool | `true` | Enable database |
| `db.driver` | `DB_DRIVER` | string | `postgres` | DB driver: `postgres`, `mariadb`, `sqlite` |
| `db.postgres.host` | `DB_POSTGRES_HOST` | string | `postgres` | PostgreSQL host |
| `db.postgres.port` | `DB_POSTGRES_PORT` | int | `5432` | PostgreSQL port |
| `db.postgres.user` | `DB_POSTGRES_USER` | string | `postgres` | PostgreSQL user |
| `db.postgres.password` | `DB_POSTGRES_PASSWORD` | string | `postgrespassword` | PostgreSQL password |
| `db.postgres.dbname` | `DB_POSTGRES_DBNAME` | string | `bkit_db` | PostgreSQL database name |
| `db.postgres.sslmode` | `DB_POSTGRES_SSLMODE` | string | `disable` | PostgreSQL SSL mode |
| `db.mariadb.host` | `DB_MARIADB_HOST` | string | `mariadb` | MariaDB host |
| `db.mariadb.port` | `DB_MARIADB_PORT` | int | `3306` | MariaDB port |
| `db.mariadb.user` | `DB_MARIADB_USER` | string | `bkit_user` | MariaDB user |
| `db.mariadb.password` | `DB_MARIADB_PASSWORD` | string | `mariadbpassword` | MariaDB password |
| `db.mariadb.dbname` | `DB_MARIADB_DBNAME` | string | `bkit_db` | MariaDB database name |
| `db.sqlite.path` | `DB_SQLITE_PATH` | string | `data.db` | SQLite file path |
| `db.max_open_conns` | `DB_MAX_OPEN_CONNS` | int | `10` | Max open DB connections |
| `db.max_idle_conns` | `DB_MAX_IDLE_CONNS` | int | `5` | Max idle DB connections |

---

## 4. Code — What You Can Change

### Entry Point — `cmd/main.go`

The entrypoint is a thin orchestrator. It:

1. Creates a root context with OS signal handling (`SIGINT` / `SIGTERM`)
2. Loads config via `config.Load(ctx)`
3. Initializes the `bsuite.Service` container (DB, KV, Telemetry)
4. Registers "runnables" (the HTTP server + background worker) with `brun.Manager`
5. Starts everything concurrently and waits for shutdown

**You should change:** Add new `brun.Runnable` implementations here if you need additional concurrent processes (e.g., a gRPC server, a queue consumer).

### Config Loader — `internal/config/config.go`

Contains the custom `EnvSource` that maps flat environment variables (`SERVER_PORT`) into the nested config tree (`server.port`). Only processes env vars whose first segment matches `server`, `telemetry`, `kv`, or `db`.

**You should change:** If you add a new top-level config section (e.g., `auth:`), add its prefix to the allowlist on line 45:

```go
if firstPart != "server" && firstPart != "telemetry" && firstPart != "kv" && firstPart != "db" {
    continue
}
// Add your new section ↑ e.g.: && firstPart != "auth"
```

### HTTP Server — `internal/server/server.go`

An Echo v4 server registered as a `brun.Runnable`. It handles:

- Middleware stack (Recover, RequestID, Logger, CORS, optional OTel tracing)
- Static asset serving (`/assets`)
- Routes: `GET /` (dashboard), `GET /api/status` (HTMX status grid), `GET /metrics` (Prometheus)
- Graceful shutdown on context cancellation

**You should change:**

- **Add routes** — define new handlers and register them (e.g., `e.GET("/api/users", s.handleUsers)`)
- **Add middleware** — insert custom auth, rate-limiting, etc.
- **Remove the dashboard** — replace the home route with your own API or pages.

### Background Worker — `internal/worker/worker.go`

A ticker-based background loop registered as a `brun.Runnable`. Runs health checks every 15 seconds and demonstrates KV writes.

**You should change:** Replace the `performCheck` logic with your own background tasks — queue processing, scheduled jobs, cache warming, etc.

### Templates — `views/pages/`

Templ templates for the built-in dashboard UI. Uses `templUI` components (auto-generated into `views/components/`).

| File | Purpose |
|---|---|
| `layout.templ` | Base HTML layout (head, fonts, scripts, CSS) |
| `home.templ` | Dashboard page with status cards, config tables, and quickstart guides |

**You should change:** Replace or extend these templates for your own UI. If building a pure API, you can remove the `views/` directory entirely and strip out the templ/templUI dependencies.

### Static Assets — `assets/`

| Path | Purpose |
|---|---|
| `assets/css/input.css` | Tailwind v4 entry file — add custom CSS imports here |
| `assets/js/lib/htmx.js` | HTMX library for partial page updates |

**You should change:** Add your own JavaScript, images, or CSS files here. They're served at `/assets/*`.

---

## 5. Running the Project

### Option 1: Full Docker Stack (Recommended)

Start everything — app, Postgres, MariaDB, Redis — with one command:

```bash
docker compose up
```

The app will be available at **http://localhost:8080** with full hot-reloading. Edit any Go, Templ, or CSS file and the server rebuilds automatically.

### Option 2: Docker Services + Local App

Start only the databases, then run the app natively:

```bash
# Start databases
docker compose up -d postgres mariadb redis

# Run locally with hot-reload
make dev
```

### Option 3: Production Build

```bash
# Build a static binary
make build

# Run it
./bin/server
```

Or build the production Docker image:

```bash
docker build -f Dockerfile -t my-app .
```

---

## 6. Taking Over for Your Real Project

Follow these steps to transform this template into your own application.

### Step 1: Rename the Module

Update the Go module path in `go.mod` from `github.com/brian-nunez/bkit-api-template` to your own:

```bash
# Find and replace across all Go files
grep -rl "github.com/brian-nunez/bkit-api-template" --include="*.go" --include="go.mod" . \
  | xargs sed -i '' 's|github.com/brian-nunez/bkit-api-template|github.com/yourorg/yourapp|g'
```

### Step 2: Update Configuration Defaults

Edit **`config.yaml`** with your values:

```yaml
server:
  port: 3000                     # Your preferred port

telemetry:
  service_name: "my-service"     # Your service name
  environment: "development"

db:
  driver: "postgres"             # Pick your database
  postgres:
    host: "my-db-host"
    user: "my_user"
    password: "my_password"
    dbname: "my_database"
```

Then mirror those values in **`compose.yml`** environment variables if running via Docker:

```yaml
environment:
  - SERVER_PORT=3000
  - TELEMETRY_SERVICE_NAME=my-service
  - DB_POSTGRES_HOST=my-db-host
  - DB_POSTGRES_USER=my_user
  - DB_POSTGRES_PASSWORD=my_password
  - DB_POSTGRES_DBNAME=my_database
```

### Step 3: Remove Unused Database Services

If you only need Postgres, remove the MariaDB and Redis services from `compose.yml`, remove their volumes, and disable them in config:

```yaml
# config.yaml
kv:
  enabled: false    # or switch driver to "local" if you still want in-memory caching

db:
  driver: "postgres"
```

Remove the corresponding `depends_on` entries and environment variables from the `app` service.

### Step 4: Add Your Routes and Business Logic

1. **Add new route handlers** in `internal/server/server.go`:

   ```go
   // In the Run method, after existing routes:
   e.GET("/api/users", s.handleListUsers)
   e.POST("/api/users", s.handleCreateUser)
   ```

2. **Create new packages** under `internal/` for your domain logic:

   ```text
   internal/
   ├── server/     # HTTP layer
   ├── worker/     # Background jobs
   ├── models/     # ← Your data models
   ├── repository/ # ← Your database queries
   └── service/    # ← Your business logic
   ```

3. **Access the database** from any handler via the service container:

   ```go
   func (s *Server) handleListUsers(c echo.Context) error {
       db := s.container.DB()
       rows, err := db.Query(c.Request().Context(), "SELECT id, name FROM users")
       // ...
   }
   ```

### Step 5: Add New Config Sections (Optional)

If you need new top-level config (e.g., `auth`):

1. Add the section to `config.yaml`:

   ```yaml
   auth:
     jwt_secret: "change-me"
     token_ttl: 3600
   ```

2. Add env var overrides to `compose.yml`:

   ```yaml
   - AUTH_JWT_SECRET=my-secret
   - AUTH_TOKEN_TTL=7200
   ```

3. Update the env var allowlist in `internal/config/config.go`:

   ```go
   if firstPart != "server" && firstPart != "telemetry" && firstPart != "kv" && firstPart != "db" && firstPart != "auth" {
       continue
   }
   ```

4. Access in code:

   ```go
   secret := cfg.String("auth.jwt_secret")
   ttl := cfg.Int("auth.token_ttl")
   ```

### Step 6: Replace or Remove the Dashboard

The template ships with a built-in dashboard at `/`. For a headless API:

1. Remove `views/pages/home.templ` and `views/pages/layout.templ`
2. Remove the template-related imports and renderer from `server.go`
3. Replace the `GET /` route with your own handler or redirect
4. Remove the `templui`, `templ`, and Tailwind dependencies from `Makefile` and `Dockerfile` if not needed

### Step 7: Update Docker Configuration

- **`Dockerfile`** — Update if you rename directories or remove dependencies
- **`Dockerfile.dev`** — Same as above
- **`compose.yml`** — Update container names, ports, and volume mounts to match your project
- **`.air.toml`** — Update `include_dir` if you add new directories that should trigger hot-reload

### Step 8: Clean Up Template Artifacts

```bash
# Remove template-specific files
rm -f ONBOARDING.md          # This file — once you've onboarded
rm -f .templui.json          # Only needed if using templUI

# Update README.md with your own project description
```

---

## Quick Reference: Config Cheat Sheet

| I want to... | Edit `config.yaml` | Edit `compose.yml` env vars |
|---|---|---|
| Change the server port | `server.port: 9090` | `SERVER_PORT=9090` |
| Switch to MariaDB | `db.driver: "mariadb"` | `DB_DRIVER=mariadb` |
| Disable telemetry | `telemetry.enabled: false` | `TELEMETRY_ENABLED=false` |
| Use local KV instead of Redis | `kv.driver: "local"` | `KV_DRIVER=local` |
| Change the Redis address | `kv.redis.addr: "my-redis:6379"` | `KV_REDIS_ADDR=my-redis:6379` |
| Update DB credentials | `db.postgres.password: "newpass"` | `DB_POSTGRES_PASSWORD=newpass` |
| Rename the service | `telemetry.service_name: "my-svc"` | `TELEMETRY_SERVICE_NAME=my-svc` |
