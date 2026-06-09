# bkit-api-template

A clean, modern, and easily expandable Go api template designed for local development and production. This project serves as a boilerplate for future applications.

It showcases full integration of the **bkit** micro-library ecosystem, **Echo v4**, **Templ**, **Tailwind CSS v4**, and **templUI** components with a unified hot-reloading development container.

ALL functionality is yours to modify. This template only glues together the libraries and provides a starting point for your application. The included features are examples of how to use the libraries together, but you can remove or replace any part as needed.

## 🚀 Features

*   **Unified Configuration Loader (`bconfig`):** Supports configuration merging from `config.yaml` defaults and environment variable overrides set in `compose.yml`. Env vars take precedence over the YAML file — you can use either or both.
*   **Database Client Integration (`bdb`):** Complete wrapper support for PostgreSQL, MariaDB (MySQL), and SQLite with connection pooling.
*   **Key-Value Cache Store (`bkv`):** Built-in support for Redis and local memory stores.
*   **Observability & Telemetry (`btelemetry`):** Tracing and metrics instrumentation using OpenTelemetry and Prometheus scraping endpoint.
*   **Service Container Wrapper (`bsuite`):** Unified initialization of config, db, kv, and telemetry.
*   **Concurrent Runnable Orchestrator (`brun`):** Manager that runs both the HTTP server and background worker threads concurrently and handles graceful shutdown.
*   **Tailwind CSS v4 & Templ UI:** Server-side rendered components utilizing `templui` and utility CSS, downloaded and built entirely on compilation.
*   **Zero-Dependency Git Repo:** All components, compiled CSS, and templates are generated during build/compilation and are excluded from version control.

---

## 🔧 Configuration

Configuration can be set in **two places** — they are merged at startup with environment variables winning over file values:

| Method | File | When to Use |
|---|---|---|
| **YAML Config** | `config.yaml` | Default values checked into the repo. Edit for local development baselines. |
| **Environment Variables** | `compose.yml` → `environment:` | Per-environment overrides (dev, staging, prod). No code changes needed. |

**How it works:** The config loader reads `config.yaml` first, then applies environment variable overrides on top. Environment variables use `SCREAMING_SNAKE_CASE` and map to the nested YAML keys by replacing dots/nesting with underscores.

```
# YAML (config.yaml)            # Environment Variable (compose.yml)
server:
  port: 8080                     SERVER_PORT=9090

db:
  driver: postgres               DB_DRIVER=mariadb
  postgres:
    host: postgres               DB_POSTGRES_HOST=my-db-host

telemetry:
  enabled: true                  TELEMETRY_ENABLED=false
```

> For the complete configuration reference and step-by-step onboarding guide, see **[ONBOARDING.md](./ONBOARDING.md)**.

---

## 📂 Project Structure

```text
├── assets/                  # Public web assets
│   ├── css/
│   │   └── input.css        # Tailwind v4 configuration file
│   └── js/
│       └── lib/
│           └── htmx.js      # HTMX framework
├── cmd/
│   └── main.go              # App entrypoint (initializes bsuite & brun manager)
├── internal/
│   ├── config/              # Configuration loader
│   ├── server/              # Echo HTTP server wrapper (brun.Runnable)
│   └── worker/              # Background worker loop (brun.Runnable)
├── views/
│   ├── pages/               # Page templates (home, layout)
│   └── utils/               # Generated templui utility functions
├── .air.toml                # Air hot-reloading tool configuration
├── .templui.json            # templUI CLI configuration file
├── compose.yml              # Multi-container local environment (Postgres, MariaDB, Redis)
├── Dockerfile               # Production multi-stage Docker build
├── Dockerfile.dev           # Development container configuration
├── Makefile                 # Build and automation commands
└── config.yaml              # Local default configurations
```

---

## ⚙️ Development Commands

We provide a comprehensive `Makefile` to automate common development workflows.

### Prerequisites

The Makefile automatically detects and installs the required Go CLIs (`templ`, `air`, `templui`) to your local `~/go/bin` if they are not already in your `$PATH`.

### Local Execution

1.  **Start Sibling Databases (Postgres, MariaDB, Redis):**
    ```bash
    docker compose up -d postgres mariadb redis
    ```
2.  **Start Hot-Reloading Development Server:**
    ```bash
    make dev
    ```
    This command will:
    *   Download and install `templ`, `air`, and `templui` tools.
    *   Pull the latest components from the `templui` CDN into `views/components/`.
    *   Watch and compile templates on change.
    *   Watch and compile CSS on change using `tailwindcss`.
    *   Run `air` to rebuild the Go binary on any source file modification.

3.  **Build Production Binary:**
    ```bash
    make build
    ```
    Creates a production-ready statically linked binary in `bin/server`.

4.  **Clean Generated Files:**
    ```bash
    make clean
    ```
    Removes all build artifacts, generated components, compiled CSS, and generated templ Go files.

---

## 🐳 Docker Deployment

The application features full Docker integration.

### Running the Entire Stack in Development (Hot-Reloading)

```bash
docker compose up
```

This starts the application along with Postgres, MariaDB, and Redis. The source folder is mounted directly into the dev container. Any changes to the Go files, templ files, or libraries will trigger hot reloading immediately.

### Production Build

To build the production docker image manually:

```bash
docker build -f Dockerfile -t bkit-api-template .
```

Or run it using docker-compose production files.
