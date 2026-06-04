# Northstar CMDB Justfile
# A collection of simple shortcut recipes to streamline local development, testing, and documentation operations.

# Default: list all available recipes
default:
    @just --list

# Start the Golang REST API backend server (port 8000)
backend:
    cd backend && go run cmd/server/main.go

# Start the Nuxt 3 frontend development dashboard (port 3000)
frontend:
    cd frontend && npm run dev

# Run Go backend unit & integration tests
test-backend:
    cd backend && go test -v ./...

# Run Vitest frontend unit tests
test-frontend:
    cd frontend && npm run test

# Run Playwright E2E browser tests against a real, isolated database
test-e2e:
    cd frontend && npx playwright test

# Install MkDocs documentation dependencies
docs-setup:
    pip install mkdocs-material

# Start the local hot-reloading developer documentation server (port 8000)
docs-serve:
    mkdocs serve

# Compile documentation markdown files into static optimized HTML pages
docs-build:
    mkdocs build
