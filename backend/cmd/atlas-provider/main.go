package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"cmdb-backend/internal/models"
)

func main() {
	stmts, err := gormschema.New("sqlite").Load(
		&models.Asset{},
		&models.MaintenanceWindow{},
		&models.PortTypeProfile{},
		&models.NetworkInterface{},
		&models.AuditLog{},
		&models.AssetRelationship{},
		&models.Manufacturer{},
		&models.Category{},
		&models.SubGroup{},
		&models.DeviceModel{},
		&models.DeviceModelRevision{},
		&models.Permission{},
		&models.Group{},
		&models.User{},
		&models.DatacenterType{},
		&models.DatacenterFloor{},
		&models.DatacenterWall{},
		&models.Datacenter{},
		&models.Rack{},
		&models.CustomField{},
		&models.Webhook{},
		&models.LicenseAgreement{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
