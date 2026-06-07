# Northstar CMDB Project - Gemini CLI Instructions

This file contains foundational mandates, architectural rules, and project context for Gemini CLI and other AI agents. Adhering to these guidelines reduces token consumption, keeps context clean, and ensures high-quality engineering.

### 🛡️ Branding & Ownership Mandate (Strict)
- The official owner of this project is **Astrona** (owned by **astrona.io**).
- The official product/brand name of the CMDB platform is **Northstar** (or **Northstar CMDB**).
- No other company names or legacy logo namespaces are permitted. All system API groups, custom YAML manifests, configuration templates, or domains MUST utilize **astrona.io** (e.g. `apiVersion: network.astrona.io/v1alpha1`).

---

## 🚀 Context & Token Efficiency Mandates

To minimize unnecessary context usage and speed up interactions, AI agents MUST follow these instructions:

### 1. Directories to Ignore
The following folders are auto-generated, heavy, or non-essential and **MUST NOT** be read, indexed, or traversed:
- `node_modules/` (Frontend dependencies)
- `.nuxt/` (Nuxt build cache)
- `.output/` (Nuxt production build)
- `.venv/` (Python virtual environment)
- `backend_python_backup/` (Legacy Python codebase)
- `.devbox/` (Nix environment files)

> **Note:** These paths are explicitly defined in the `.geminiignore` file at the root.

### 2. Surgical Reading & Writing
- **Do not read whole files** unless they are small (< 100 lines). Use the targeted line-range parameters `start_line` and `end_line` with `read_file` to only fetch relevant blocks.
- **Do not use `cat` or `grep` inside bash commands.** Utilize specific tools like `grep_search` and `read_file` which are far more token-efficient and safe.
- **Prefer `replace` over `write_file`** for existing files to avoid rewriting unchanged code and blowing up the history.

---

## 🏗️ Architecture & Conventions

### Backend (Golang + GORM + Echo)
- Located in `backend/`
- Enforces standard Go layout:
  - `cmd/server/main.go` — Entrypoint, Echo router, middlewares.
  - `internal/database/` — GORM SQLite setup & migrations.
  - `internal/models/` — DB GORM structs. Has computed property calculation logic.
  - `internal/handlers/` — HTTP handlers.
- **Properties Schema:** The `properties` field in models is validated using structs defined in `internal/handlers/properties_schema.go` before updates.

### Frontend (Nuxt 3 + Nuxt UI)
- Located in `frontend/`
- **Composables:** API interactions must go through composables inside `frontend/composables/` (e.g., `useAssets.ts`, `useMaintenance.ts`) rather than direct `$fetch` or `useFetch`.
- **Componentization:** Keep pages clean. Extract tables, modals, and property displays into single-responsibility components in `frontend/components/`.

---

## 🧪 Testing Requirements

- **Backend:** Tests use `go test -v ./...` under `backend/`. Unit tests live alongside models, integration tests under `internal/handlers/assets_test.go` (running on an in-memory SQLite).
- **Frontend Unit Tests:** Setup with Vitest. Run using `npx vitest run`.
- **E2E / UX Tests:** Setup with Playwright under `frontend/e2e/`. Run using `npx playwright test`.

---

## 🗄️ Database Migrations & Versioning Mandate (Strict)

Database schemas are strictly managed and tracked using **Goose** versioned SQL migrations compiled directly into the binary using Go's `//go:embed` file system under `backend/internal/database/migrations/`.

To prevent database conflicts across deployments and releases, AI agents and developers **MUST** adhere to the following pipeline whenever modifying schemas or database seeding defaults:

### 1. Release Version Verification
- Before creating a new migration file, you **MUST** inspect the latest published releases on GitHub:
  ➡️ **`https://github.com/Astrona-io/northstar/releases`**
- Identify the current active release version. The next migration file **MUST** be named using the next sequence integer prefix and the next matched semantic pattern:
  **`vX.(new-db-version).0`** as a suffix.

### 2. Hybrid Filename Format Enforce
Goose requires sequential files to start with a valid integer, so you must use this hybrid naming scheme:
- `0001_v0.0.1_initial_schema.sql` (Baseline GORM tables)
- `0002_v0.0.2_seed_data.sql` (Fallback categories, webhooks, and types)
- `0003_v0.0.3_your_change_summary.sql` (Next release target)

### 3. CI/CD Migration Testing
- All unit and integration test setups (`setupTestDB()`) **MUST** execute the Goose migration ledger natively (`database.RunMigrations(db)`) in memory on every run.
- This guarantees that any SQL schema alterations, missing tables, or foreign-key violations are automatically caught and flagged during testing before merging pull requests.

### 4. Modular & Table-Specific Migration Chunks
- Do not create single, massive migration files for large features or release stages. Instead, you **MUST** split migrations into smaller, highly readable, modular "chunks" (e.g. separating table schemas, indices, database seed data insertions, or subsequent column alterations into distinct, separate sequential files).
- **Table-Specific Naming Enforced**: Each migration file dealing with schema creations or modifications **MUST** include the exact target table name (or join-table name) inside the filename (e.g. `0003_v0.0.3_add_is_imported_to_device_models.sql` or `0004_v0.0.4_create_users_table.sql`). This ensures perfect database auditability, keeps files self-contained, and allows developers to easily trace the history of a single database table on disk.
