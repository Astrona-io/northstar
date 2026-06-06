package router

import (
	"net/http"

	"cmdb-backend/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes configures all standard REST API endpoints, security middleware,
// and CORS controls on the provided Echo server instance.
func RegisterRoutes(e *echo.Echo) {
	// Logger & Recover Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	// API Routing
	// Registering both slash and non-slash endpoints for bulletproof frontend compatibility.
	e.GET("/api/assets", handlers.ReadAssets)
	e.GET("/api/assets/", handlers.ReadAssets)
	e.GET("/api/assets/:asset_id", handlers.ReadAsset)
	e.GET("/api/assets/:asset_id/", handlers.ReadAsset)

	// Authentication Routes (Phase 2)
	e.POST("/api/auth/login", handlers.Login)
	e.POST("/api/auth/login/", handlers.Login)

	// Fine-Grained RBAC Permission Middlewares (Phase 1 Advanced RBAC)
	pAssetWrite := handlers.RBACMiddleware("asset:write")
	pAssetDelete := handlers.RBACMiddleware("asset:delete")
	pCatalogWrite := handlers.RBACMiddleware("catalog:write")

	e.POST("/api/assets", handlers.CreateAsset, pAssetWrite)
	e.POST("/api/assets/", handlers.CreateAsset, pAssetWrite)
	e.PUT("/api/assets/:asset_id", handlers.UpdateAsset, pAssetWrite)
	e.PUT("/api/assets/:asset_id/", handlers.UpdateAsset, pAssetWrite)

	e.DELETE("/api/assets/:asset_id", handlers.DeleteAsset, pAssetDelete)
	e.DELETE("/api/assets/:asset_id/", handlers.DeleteAsset, pAssetDelete)

	// SBOM SPDX Upload endpoint (Phase 6 Zero Mocks)
	e.POST("/api/assets/:id/sbom", handlers.ParseSBOM, pAssetWrite)
	e.POST("/api/assets/:id/sbom/", handlers.ParseSBOM, pAssetWrite)

	// Maintenance Window Routes
	e.GET("/api/assets/:asset_id/maintenance", handlers.ReadMaintenanceWindows)
	e.GET("/api/assets/:asset_id/maintenance/", handlers.ReadMaintenanceWindows)

	e.POST("/api/assets/:asset_id/maintenance", handlers.CreateMaintenanceWindow, pAssetWrite)
	e.POST("/api/assets/:asset_id/maintenance/", handlers.CreateMaintenanceWindow, pAssetWrite)

	e.PUT("/api/maintenance/:window_id", handlers.UpdateMaintenanceWindow, pAssetWrite)
	e.PUT("/api/maintenance/:window_id/", handlers.UpdateMaintenanceWindow, pAssetWrite)

	e.DELETE("/api/maintenance/:window_id", handlers.DeleteMaintenanceWindow, pAssetDelete)
	e.DELETE("/api/maintenance/:window_id/", handlers.DeleteMaintenanceWindow, pAssetDelete)

	// Audit Log Routes (Phase 2)
	e.GET("/api/assets/:asset_id/logs", handlers.ReadAuditLogs)
	e.GET("/api/assets/:asset_id/logs/", handlers.ReadAuditLogs)

	// Asset Relationships Routes (Phase 8/Topology Links)
	e.GET("/api/assets/:asset_id/relationships", handlers.ReadRelationships)
	e.GET("/api/assets/:asset_id/relationships/", handlers.ReadRelationships)
	e.POST("/api/assets/:asset_id/relationships", handlers.SyncRelationships, pAssetWrite)
	e.POST("/api/assets/:asset_id/relationships/", handlers.SyncRelationships, pAssetWrite)

	// Device Catalog Routes (Plan - Manufacturer & Models)
	e.GET("/api/devices", handlers.ReadDeviceModels)
	e.GET("/api/devices/", handlers.ReadDeviceModels)
	e.POST("/api/devices", handlers.CreateDeviceModel, pCatalogWrite)
	e.POST("/api/devices/", handlers.CreateDeviceModel, pCatalogWrite)
	e.GET("/api/devices/:device_id", handlers.ReadDeviceModel)
	e.GET("/api/devices/:device_id/", handlers.ReadDeviceModel)
	e.PUT("/api/devices/:device_id", handlers.UpdateDeviceModel, pCatalogWrite)
	e.PUT("/api/devices/:device_id/", handlers.UpdateDeviceModel, pCatalogWrite)
	e.DELETE("/api/devices/:device_id", handlers.DeleteDeviceModel, pCatalogWrite)
	e.DELETE("/api/devices/:device_id/", handlers.DeleteDeviceModel, pCatalogWrite)

	// Manufacturers & Categories auxiliary routes (Phase 2 Catalog Update)
	e.GET("/api/manufacturers", handlers.ReadManufacturers)
	e.GET("/api/manufacturers/", handlers.ReadManufacturers)
	e.POST("/api/manufacturers", handlers.CreateManufacturer, pCatalogWrite)
	e.POST("/api/manufacturers/", handlers.CreateManufacturer, pCatalogWrite)
	e.GET("/api/categories", handlers.ReadCategories)
	e.GET("/api/categories/", handlers.ReadCategories)

	// Auto-Discovery Routes (Phase 1)
	e.GET("/api/discovery/scan", handlers.ScanSubnet)
	e.GET("/api/discovery/scan/", handlers.ScanSubnet)
	e.POST("/api/discovery/import", handlers.ImportDevice, pAssetWrite)
	e.POST("/api/discovery/import/", handlers.ImportDevice, pAssetWrite)
	e.POST("/api/discovery/cloud-sync", handlers.CloudSync, pAssetWrite)
	e.POST("/api/discovery/cloud-sync/", handlers.CloudSync, pAssetWrite)

	// Webhooks Integration Routes (Phase 3 Webhooks)
	e.GET("/api/webhooks", handlers.ReadWebhooks, pAssetWrite)
	e.GET("/api/webhooks/", handlers.ReadWebhooks, pAssetWrite)
	e.POST("/api/webhooks", handlers.CreateWebhook, pAssetWrite)
	e.POST("/api/webhooks/", handlers.CreateWebhook, pAssetWrite)
	e.DELETE("/api/webhooks/:id", handlers.DeleteWebhook, pAssetDelete)
	e.DELETE("/api/webhooks/:id/", handlers.DeleteWebhook, pAssetDelete)

	// Global Audit Logs (Phase 6 Zero Mocks)
	e.GET("/api/audit-logs", handlers.ReadGlobalAuditLogs, pAssetWrite)
	e.GET("/api/audit-logs/", handlers.ReadGlobalAuditLogs, pAssetWrite)

	// Custom Fields Routes (Phase 4 Custom Fields)
	e.GET("/api/custom-fields", handlers.ReadCustomFields)
	e.GET("/api/custom-fields/", handlers.ReadCustomFields)
	e.POST("/api/custom-fields", handlers.CreateCustomField, pAssetWrite)
	e.POST("/api/custom-fields/", handlers.CreateCustomField, pAssetWrite)
	e.DELETE("/api/custom-fields/:id", handlers.DeleteCustomField, pAssetDelete)
	e.DELETE("/api/custom-fields/:id/", handlers.DeleteCustomField, pAssetDelete)

	// Dynamic Category Routes (Phase 1 Admin Category Manager)
	e.GET("/api/categories", handlers.ReadCategories)
	e.GET("/api/categories/", handlers.ReadCategories)
	e.POST("/api/categories", handlers.CreateCategory, pAssetWrite)
	e.POST("/api/categories/", handlers.CreateCategory, pAssetWrite)
	e.DELETE("/api/categories/:id", handlers.DeleteCategory, pAssetDelete)
	e.DELETE("/api/categories/:id/", handlers.DeleteCategory, pAssetDelete)
	e.POST("/api/categories/sub-groups", handlers.CreateSubGroup, pAssetWrite)
	e.POST("/api/categories/sub-groups/", handlers.CreateSubGroup, pAssetWrite)

	// Live Network Latency Monitoring Route (Phase 1 Latency)
	e.GET("/api/monitoring/ping", handlers.GetLatencyPing)
	e.GET("/api/monitoring/ping/", handlers.GetLatencyPing)

	// Administrative License Agreement Routes (AGPLv3 Compliance)
	e.GET("/api/license/status", handlers.GetLicenseStatus, pAssetWrite)
	e.GET("/api/license/status/", handlers.GetLicenseStatus, pAssetWrite)
	e.POST("/api/license/accept", handlers.AcceptLicense, pAssetWrite)
	e.POST("/api/license/accept/", handlers.AcceptLicense, pAssetWrite)

	// Datacenter & Rack Unit DCIM Routes (Phase 6 DCIM Mapping)
	e.GET("/api/datacenter-types", handlers.ReadDatacenterTypes)
	e.GET("/api/datacenter-types/", handlers.ReadDatacenterTypes)
	e.POST("/api/datacenter-types", handlers.CreateDatacenterType, pAssetWrite)
	e.POST("/api/datacenter-types/", handlers.CreateDatacenterType, pAssetWrite)

	e.GET("/api/datacenters", handlers.ReadDatacenters)
	e.GET("/api/datacenters/", handlers.ReadDatacenters)
	e.POST("/api/datacenters", handlers.CreateDatacenter, pAssetWrite)
	e.POST("/api/datacenters/", handlers.CreateDatacenter, pAssetWrite)
	e.PUT("/api/datacenters/:id", handlers.UpdateDatacenter, pAssetWrite)
	e.PUT("/api/datacenters/:id/", handlers.UpdateDatacenter, pAssetWrite)
	e.DELETE("/api/datacenters/:id", handlers.DeleteDatacenter, pAssetDelete)
	e.DELETE("/api/datacenters/:id/", handlers.DeleteDatacenter, pAssetDelete)

	e.POST("/api/datacenter-floors", handlers.CreateDatacenterFloor, pAssetWrite)
	e.POST("/api/datacenter-floors/", handlers.CreateDatacenterFloor, pAssetWrite)
	e.PUT("/api/datacenter-floors/:id", handlers.UpdateDatacenterFloor, pAssetWrite)
	e.PUT("/api/datacenter-floors/:id/", handlers.UpdateDatacenterFloor, pAssetWrite)
	e.DELETE("/api/datacenter-floors/:id", handlers.DeleteDatacenterFloor, pAssetDelete)
	e.DELETE("/api/datacenter-floors/:id/", handlers.DeleteDatacenterFloor, pAssetDelete)

	e.POST("/api/dcim/walls", handlers.CreateDatacenterWall, pAssetWrite)
	e.POST("/api/dcim/walls/", handlers.CreateDatacenterWall, pAssetWrite)
	e.DELETE("/api/dcim/walls/:id", handlers.DeleteDatacenterWall, pAssetDelete)
	e.DELETE("/api/dcim/walls/:id/", handlers.DeleteDatacenterWall, pAssetDelete)

	e.GET("/api/racks", handlers.ReadRacks)
	e.GET("/api/racks/", handlers.ReadRacks)
	e.POST("/api/racks", handlers.CreateRack, pAssetWrite)
	e.POST("/api/racks/", handlers.CreateRack, pAssetWrite)
	e.DELETE("/api/racks/:id", handlers.DeleteRack, pAssetDelete)
	e.DELETE("/api/racks/:id/", handlers.DeleteRack, pAssetDelete)
	e.GET("/api/racks/:id", handlers.ReadRack)
	e.GET("/api/racks/:id/", handlers.ReadRack)
	e.PUT("/api/racks/:id/coordinates", handlers.UpdateRackCoordinates, pAssetWrite)
	e.PUT("/api/racks/:id/coordinates/", handlers.UpdateRackCoordinates, pAssetWrite)

	// User Management Routes (Phase 1 Admin User overrides)
	e.GET("/api/users", handlers.ReadUsers, pAssetWrite)
	e.GET("/api/users/", handlers.ReadUsers, pAssetWrite)
	e.PUT("/api/users/:id", handlers.UpdateUser, pAssetDelete)
	e.PUT("/api/users/:id/", handlers.UpdateUser, pAssetDelete)
	e.GET("/api/groups", handlers.ReadGroups, pAssetWrite)
	e.GET("/api/groups/", handlers.ReadGroups, pAssetWrite)
	e.GET("/api/permissions", handlers.ReadPermissions, pAssetWrite)
	e.GET("/api/permissions/", handlers.ReadPermissions, pAssetWrite)

	// Network Interface Routes (Phase 7)
	e.GET("/api/assets/:asset_id/interfaces", handlers.ReadNetworkInterfaces)
	e.GET("/api/assets/:asset_id/interfaces/", handlers.ReadNetworkInterfaces)

	// Create physical/virtual port connections on target NICs
	e.POST("/api/assets/:asset_id/interfaces", handlers.CreateNetworkInterface, pAssetWrite)
	e.POST("/api/assets/:asset_id/interfaces/", handlers.CreateNetworkInterface, pAssetWrite)

	// Update existing network interface card configuration details
	e.PUT("/api/interfaces/:interface_id", handlers.UpdateNetworkInterface, pAssetWrite)
	e.PUT("/api/interfaces/:interface_id/", handlers.UpdateNetworkInterface, pAssetWrite)

	// Delete specific physical/virtual port interfaces safely
	e.DELETE("/api/interfaces/:interface_id", handlers.DeleteNetworkInterface, pAssetDelete)
	e.DELETE("/api/interfaces/:interface_id/", handlers.DeleteNetworkInterface, pAssetDelete)
}
