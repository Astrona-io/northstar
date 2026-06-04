package database

import (
	"log"
	"time"

	"cmdb-backend/internal/models"

	"gorm.io/gorm"
)

type MigrationFunc func(db *gorm.DB) error

// versionedMigrations holds all custom declarative migration steps (similar to alembic).
var versionedMigrations = map[int]struct {
	Description string
	Migrate     MigrationFunc
}{
	1: {
		Description: "Base CMDB schema auto-migration and initial seed data",
		Migrate: func(db *gorm.DB) error {
			// Baseline GORM migrations are handled on startup
			return nil
		},
	},
	2: {
		Description: "Seed standard webhooks and custom fields defaults (Phase 1/4)",
		Migrate: func(db *gorm.DB) error {
			// Seed a default Custom Field to show immediately
			field := models.CustomField{
				Name:       "warranty_expiry",
				FieldType:  "text",
				AssetType:  "Server",
				IsRequired: false,
			}
			db.FirstOrCreate(&field, models.CustomField{Name: field.Name})

			// Seed a default webhook
			webhook := models.Webhook{
				URL:   "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
				Event: "asset:create",
			}
			db.FirstOrCreate(&webhook, models.Webhook{URL: webhook.URL, Event: webhook.Event})

			log.Println("[Migration Engine] Completed Version 2 Schema Update successfully.")
			return nil
		},
	},
	3: {
		Description: "Seed standard Datacenter Types (Phase 3 DCIM Upgrade)",
		Migrate: func(db *gorm.DB) error {
			types := []string{"on-prem", "homelab", "air-gap", "cloud:aws", "cloud:gcp", "cloud:azure"}
			for _, t := range types {
				var dcType models.DatacenterType
				db.FirstOrCreate(&dcType, models.DatacenterType{Name: t})
			}

			// Map existing Datacenters to use the Type IDs instead of raw strings
			var dcs []models.Datacenter
			db.Find(&dcs)
			for _, dc := range dcs {
				if dc.Type != "" && dc.DatacenterTypeID == nil {
					var matchedType models.DatacenterType
					if err := db.Where("name = ?", dc.Type).First(&matchedType).Error; err == nil {
						dc.DatacenterTypeID = &matchedType.ID
						db.Save(&dc)
					}
				}
			}

			log.Println("[Migration Engine] Completed Version 3 Schema Update successfully.")
			return nil
		},
	},
	4: {
		Description: "Seed standard asset categories and subtype sub-groups (Phase 1 Dynamic Categories - Reconciled via bootstrap.yaml)",
		Migrate: func(db *gorm.DB) error {
			// Deprecated: Hardcoded seeding loops are now managed declaratively via .astrona-bootstrap/bootstrap.yaml!
			log.Println("[Migration Engine] Completed Version 4 Schema Update successfully (Managed declaratively via bootstrap.yaml).")
			return nil
		},
	},
}

// RunMigrations executes versioned database schema migration steps. (Phase 1)
func RunMigrations(db *gorm.DB) error {
	log.Println("[Migration Engine] Checking database migration ledger status...")

	// 1. Ensure schema migrations table is created
	err := db.AutoMigrate(&models.SchemaMigration{})
	if err != nil {
		return err
	}

	// 2. Fetch completed migrations
	var applied []models.SchemaMigration
	db.Order("version asc").Find(&applied)

	appliedMap := make(map[int]bool)
	for _, m := range applied {
		appliedMap[m.Version] = true
	}

	// 3. Apply missing versions in order
	currentVersion := 0
	for version, migration := range versionedMigrations {
		if !appliedMap[version] {
			log.Printf("[Migration Engine] Applying Version %d: %s...", version, migration.Description)

			// Execute migration function
			if err := migration.Migrate(db); err != nil {
				log.Printf("[Migration Engine] Error applying migration Version %d: %v", version, err)
				return err
			}

			// Log applied version to ledger
			record := models.SchemaMigration{
				Version:   version,
				AppliedAt: time.Now(),
			}
			if err := db.Create(&record).Error; err != nil {
				return err
			}
			log.Printf("[Migration Engine] Version %d applied successfully.", version)
		}
		if version > currentVersion {
			currentVersion = version
		}
	}

	log.Printf("[Migration Engine] Database migrations up to date. Active version: v%d.", currentVersion)
	return nil
}
