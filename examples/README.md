# Northstar CMDB Container Quickstart Guide

This directory contains pre-configured orchestrations designed to launch the complete **Northstar CMDB** stack locally in seconds. It pulls public, pre-built multi-architecture container images (`amd64` and `arm64`) directly from the GitHub Container Registry (`ghcr.io/astrona-io/`).

No local compilers, Go installations, Node.js packages, or local build environments are required.

---

## 🚀 Instant SQLite Local Deployment (`sqlite-local/`)

### Prerequisites
Make sure you have **Docker** (and Docker Compose) or **Podman** (and Podman Compose) installed and running on your system.

### Step 1: Start the services
From the **`examples/sqlite-local/`** directory, run the following command to start both the Go backend and Nuxt frontend in detached (background) mode:

```bash
docker compose up -d
```

### Step 2: Access the Application
Once the containers have loaded, open your web browser and navigate to:
* **Interactive UI Frontend:** `http://localhost:3000`
* **Echo JSON REST API:** `http://localhost:8000/api`

### Step 3: Stop the services and Purge Local Data
To stop the services:
```bash
docker compose down
```
To completely wipe all local SQLite database data and start fresh on the next startup:
```bash
rm -rf data/
```

---

## 🐘 PostgreSQL & pgAdmin Local Environment (`postgres-local/`)

For advanced local evaluation, testing with a multi-user relational database, or validating database migrations with PostgreSQL, a dedicated compose file is provided inside **`examples/postgres-local/`**.

### Folder Structure & Bind Mounts
Unlike named Docker volumes, this environment mounts all data directly into a **local relative folder structure on your host machine** under `examples/postgres-local/data/`:
* **`postgres-local/data/postgres/`** ➡️ Stores all active PostgreSQL database files.
* **`postgres-local/data/pgadmin/`** ➡️ Stores your local pgAdmin settings and queries history.

This makes it incredibly easy to inspect data on disk, perform manual backups, or **instantly reset your entire testing database** simply by purging the local `./data/` folder!

### Step 1: Start the services
From the `examples/postgres-local/` directory, run:
```bash
docker compose up -d
```

### Step 2: Access the local services
* **Interactive UI Frontend:** `http://localhost:3000`
* **Echo JSON REST API:** `http://localhost:8000/api`
* **pgAdmin 4 Administration Dashboard:** `http://localhost:5050`
  * **Login Email:** `admin@astrona.io`
  * **Login Password:** `admin_secure_pass`
  * *(To connect pgAdmin to the database, register a new server inside pgAdmin and use Host: `northstar-postgres-local`, Username: `northstar`, Password: `northstar_secure_pass`, Database: `northstar`).*

### Step 3: Stop the services and Purge Local Data
To stop the services:
```bash
docker compose down
```
To completely wipe all local PostgreSQL data and start fresh on the next startup:
```bash
rm -rf data/
```

---

## 🔭 OpenTelemetry & CNCF Perses Observability Sandbox (`otel-perses/`)

For an enterprise-grade cloud-native monitoring demonstration, a pre-integrated **OpenTelemetry (OTel) + Prometheus + CNCF Perses** dashboard sandbox is provided inside **`examples/otel-perses/`**.

### The Observability Architecture
* **Northstar Backend**: Connects directly to GORM and exports core database queries, latencies, and HTTP metrics over high-speed gRPC.
* **OTel Collector**: Listens on port `4317` (OTLP gRPC) to receive telemetry from the backend, batches metrics, and exposes them on a Prometheus endpoint on port `8889`.
* **Prometheus**: Scrapes and stores metrics from the OTel Collector on port `9090`.
* **CNCF Perses**: A highly modern, GitOps-ready CNCF dashboarding engine (port `8080`) that connects directly to Prometheus to visualize live network metrics and system latencies dynamically!

### Step 1: Start the services
From the `examples/otel-perses/` directory, run:
```bash
docker compose up -d
```

### Step 2: Access the local services
* **Interactive UI Frontend:** `http://localhost:3000`
* **pgAdmin 4 Administration:** `http://localhost:5050`
* **Prometheus Explorer Panel:** `http://localhost:9090`
* **CNCF Perses Dashboard Console:** `http://localhost:8080` (Browse metrics dynamically!)

---

## 💾 Data Persistence

By default, all environment orchestrations utilize **local relative bind mounts** (e.g. `./data/` directories) rather than named Docker volumes. This ensures that your database records, configurations, and dashboards **survive container stops and restarts** while remaining 100% visible on your host disk.

If you wish to completely wipe the local database state and restart with a clean declarative bootstrap seed, stop the services and purge the `data/` folder inside the environment sub-directory:

```bash
# Stop the containers
docker compose down

# Purge local persistent data folders
rm -rf data/
```

---

## 🔍 Inspecting and Troubleshooting Logs

If you encounter any connectivity issues, check the running container logs:

```bash
# View backend log output
docker compose logs northstar-backend

# View frontend log output
docker compose logs northstar-frontend

# Follow real-time unified logs
docker compose logs -f
```
