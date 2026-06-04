# Northstar CMDB Container Quickstart Guide

This directory contains a pre-configured `docker-compose.yml` file designed to launch the complete **Northstar CMDB** stack locally in seconds. It pulls public, pre-built multi-architecture container images (`amd64` and `arm64`) directly from the GitHub Container Registry (`ghcr.io/astrona-io/`).

No local compilers, Go installations, Node.js packages, or local build environments are required.

---

## 🚀 Instant Local Deployment (30 Seconds)

### Prerequisites
Make sure you have **Docker** (and Docker Compose) or **Podman** (and Podman Compose) installed and running on your system.

### Step 1: Start the services
From this directory (`examples/`), run the following command to start both the Go backend and Nuxt frontend in detached (background) mode:

```bash
# Using Docker
docker compose up -d

# Using Podman
podman-compose up -d
```

### Step 2: Access the Application
Once the containers have loaded, open your web browser and navigate to:
* **Interactive UI Frontend:** `http://localhost:3000`
* **Echo JSON REST API:** `http://localhost:8000/api`

### Step 3: Stop the services
To stop the services and clean up the network resources without losing your persistent SQLite data, execute:

```bash
# Using Docker
docker compose down

# Using Podman
podman-compose down
```

---

## 💾 Data Persistence
By default, the backend's SQLite database file is stored in a named Docker volume (`backend-data`). This ensures that your registered datacenters, cabinets, assets, and user agreement logs **survive container stops and restarts**. 

If you wish to completely wipe the local database and restart with a clean declarative bootstrap seed, run:

```bash
# Wipe persistent volumes during shutdown
docker compose down -v
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
