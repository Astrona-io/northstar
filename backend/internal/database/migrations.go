package database

import (
	"embed"
	"log"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// RunMigrations executes versioned database schema migration steps using Goose.
func RunMigrations(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	log.Println("[Goose Engine] Running versioned schema and data migrations...")
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return err
	}

	log.Println("[Goose Engine] All database migrations applied successfully!")
	return nil
}
