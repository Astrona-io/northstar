package main

import (
	"context"
	"log"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/router"
	"cmdb-backend/internal/telemetry"
	"github.com/labstack/echo/v4"
)

// main is the primary runtime entrypoint for the Northstar CMDB Go REST API.
// It initializes the database schema, runs K8s-style bootstrapping, registers endpoints, and boots Echo.
func main() {
	// Initialize OpenTelemetry SDK (Phase 9 APM Tracing)
	ctx := context.Background()
	shutdown, err := telemetry.InitTelemetry(ctx)
	if err != nil {
		log.Printf("[OTel Engine] Failed to initialize OpenTelemetry SDK: %v", err)
	} else {
		defer shutdown()
	}

	// Initialize relational SQLite GORM database structure
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Log successful connection details for auditing
	sqlDB, err := db.DB()
	if err == nil {
		log.Printf("Successfully connected to SQLite database at: %v", sqlDB)
	}

	// Instantiate the core Echo HTTP router
	e := echo.New()

	// Enable standard OpenTelemetry Echo request tracing middleware (Phase 9)
	e.Use(telemetry.Middleware())

	// Register all REST controllers, advanced RBAC, and CORS middleware configurations (Phase 11 Refactor)
	router.RegisterRoutes(e)

	// Bind the HTTP service and listen on port 8000
	log.Println("Starting Go server on port 8000...")
	if err := e.Start(":8000"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
