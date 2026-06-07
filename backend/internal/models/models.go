package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONMap is a custom map type that implements sql.Scanner and driver.Valuer interfaces
// to allow seamless serialization of dynamic JSON properties inside SQLite.
type JSONMap map[string]any

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := json.Marshal(m)
	return string(ba), err
}

func (m *JSONMap) Scan(val any) error {
	if val == nil {
		*m = make(JSONMap)
		return nil
	}
	var ba []byte
	switch s := val.(type) {
	case []byte:
		ba = s
	case string:
		ba = []byte(s)
	default:
		return errors.New("unsupported scan type for JSONMap")
	}
	var parsed JSONMap
	err := json.Unmarshal(ba, &parsed)
	if err != nil {
		return err
	}
	*m = parsed
	return nil
}

// DatacenterType represents user-defined facility types. (Phase 3 DCIM Upgrade)
type DatacenterType struct {
	ID   string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
}

func (m *DatacenterType) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// DatacenterFloor represents a specific floor level inside a datacenter. (Phase 2 Multi-Floor)
type DatacenterFloor struct {
	ID           string           `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DatacenterID string           `gorm:"index" json:"datacenter_id"`
	Name         string           `json:"name"` // e.g., "Floor 1", "Second Floor"
	Level        int              `json:"level"`
	Width        float64          `json:"width"` // Grid boundary
	Depth        float64          `json:"depth"` // Grid boundary
	Racks        []Rack           `gorm:"foreignKey:DatacenterFloorID;constraint:OnDelete:CASCADE" json:"racks,omitempty"`
	Walls        []DatacenterWall `gorm:"foreignKey:DatacenterFloorID;constraint:OnDelete:CASCADE" json:"walls,omitempty"` // Room Wall Blueprint segments (DCIM CAD L4)
}

func (m *DatacenterFloor) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// DatacenterWall represents a drawn room barrier or wall segment inside a floor blueprint. (DCIM CAD L4)
type DatacenterWall struct {
	ID                string  `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DatacenterFloorID string  `gorm:"index" json:"datacenter_floor_id"`
	X1                float64 `json:"x1"`
	Y1                float64 `json:"y1"`
	X2                float64 `json:"x2"`
	Y2                float64 `json:"y2"`
	Thickness         float64 `gorm:"default:8" json:"thickness"` // Thickness/stroke-width of the element
	Type              string  `gorm:"default:'wall'" json:"type"` // Element type: 'wall', 'door', 'window'
	Flipped           bool    `gorm:"default:false" json:"flipped"` // Reverses the opening/swing direction
}

func (m *DatacenterWall) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Datacenter represents a physical or virtual datacenter housing racks or cloud assets. (Phase 6 DCIM)
type Datacenter struct {
	ID               string            `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name             string            `gorm:"uniqueIndex" json:"name"`
	DatacenterTypeID *string           `gorm:"index" json:"datacenter_type_id"`
	DatacenterType   *DatacenterType   `gorm:"foreignKey:DatacenterTypeID" json:"datacenter_type,omitempty"`
	Type             string            `json:"type"` // Legacy fallback
	Country          string            `json:"country"`
	City             string            `json:"city"`
	Properties       JSONMap           `gorm:"type:text" json:"properties"`
	Floors           []DatacenterFloor `gorm:"foreignKey:DatacenterID;constraint:OnDelete:CASCADE" json:"floors,omitempty"`
	Racks            []Rack            `gorm:"foreignKey:DatacenterID;constraint:OnDelete:CASCADE" json:"racks,omitempty"`
}

func (m *Datacenter) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	// Automatically create a default Floor 0 if the floors list is empty
	if len(m.Floors) == 0 {
		m.Floors = []DatacenterFloor{
			{
				ID:           uuid.New().String(),
				DatacenterID: m.ID,
				Name:         "Floor 0",
				Level:        0,
				Width:        800,
				Depth:        1200,
			},
		}
	}
	return
}

// Rack represents a physical server cabinet or virtual mounting array. (Phase 6 DCIM)
type Rack struct {
	ID                string           `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DatacenterID      string           `gorm:"index" json:"datacenter_id"`
	Datacenter        *Datacenter      `gorm:"foreignKey:DatacenterID" json:"datacenter,omitempty"`
	DatacenterFloorID *string          `gorm:"index" json:"datacenter_floor_id"`
	DatacenterFloor   *DatacenterFloor `gorm:"foreignKey:DatacenterFloorID" json:"datacenter_floor,omitempty"`
	Name              string           `json:"name"`
	HeightU           int              `gorm:"default:42" json:"height_u"`
	PlacementZone     string           `json:"placement_zone"` // e.g., "Row A", "Aisle 3"
	X                 float64          `json:"x"`              // CAD coordinate
	Y                 float64          `json:"y"`              // CAD coordinate
	Assets            []Asset          `gorm:"foreignKey:RackID" json:"assets,omitempty"`
}

func (m *Rack) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Asset represents a physical, virtual, or containerized configuration item in the CMDB. (DCIM & Containers Update)
type Asset struct {
	ID                 string              `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name               string              `gorm:"index" json:"name"`
	Type               string              `gorm:"index" json:"type"` // Server, Database, Router, Switch, Container, etc.
	Status             string              `gorm:"default:active" json:"status"`
	IPAddress          *string             `gorm:"index" json:"ip_address"`
	Description        *string             `json:"description"`
	Properties         JSONMap             `gorm:"type:text" json:"properties"`
	MaintenanceWindows []MaintenanceWindow `gorm:"foreignKey:AssetID;constraint:OnDelete:CASCADE" json:"-"`
	NetworkInterfaces  []NetworkInterface  `gorm:"foreignKey:AssetID;constraint:OnDelete:CASCADE" json:"interfaces,omitempty"`
	MaintenanceStatus  string              `gorm:"-" json:"maintenance_status"` // Calculated field (ignored by GORM)

	// DCIM Mapping Fields (Phase 6)
	RackID        *string `gorm:"index" json:"rack_id"`
	Rack          *Rack   `gorm:"foreignKey:RackID" json:"rack,omitempty"`
	RackPositionU *int    `json:"rack_position_u"` // e.g. U12 starting position height inside cabinet

	// Device Model Catalog Reference (Standard equipment ports mapping)
	DeviceModelID       *string              `gorm:"index" json:"device_model_id"`
	DeviceModel         *DeviceModel         `gorm:"foreignKey:DeviceModelID" json:"device_model,omitempty"`
	DeviceModelRevision *int                 `json:"device_model_revision"`
	DeviceModelRevisionDetails *DeviceModelRevision `gorm:"-" json:"device_model_revision_details,omitempty"`

	// Host-Guest Dependency Mapping / Containerization Fields (Phase 2)
	HostAssetID          *string `gorm:"index" json:"host_asset_id"` // Node/Server hosting this VM/Container
	HostAsset            *Asset  `gorm:"foreignKey:HostAssetID" json:"host_asset,omitempty"`
	ContainerImage       string  `json:"container_image,omitempty"`        // e.g. postgres:15-alpine (Container type only)
	ContainerPortMapping string  `json:"container_port_mapping,omitempty"` // e.g. 5432:5432 (Container type only)
}

func (m *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// MaintenanceWindow represents a scheduled maintenance event for a given Asset.
type MaintenanceWindow struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	AssetID     string    `gorm:"index" json:"asset_id"`
	Title       string    `gorm:"index" json:"title"`
	Description *string   `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `gorm:"default:scheduled" json:"status"` // scheduled, completed, cancelled
}

func (m *MaintenanceWindow) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// PortTypeProfile represents a physical port type (e.g. RJ45, SFP+) and its supported speeds.
type PortTypeProfile struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Type      string    `gorm:"uniqueIndex" json:"type"` // e.g. "RJ45", "SFP", "SFP+"
	Name      string    `json:"name"`                    // e.g. "RJ45 Copper", "SFP+ Optic"
	Speeds    string    `gorm:"type:text" json:"speeds"` // JSON array of strings, e.g. ["10 Gbps", "1 Gbps"]
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (m *PortTypeProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// NetworkInterface represents a physical or virtual network interface card (NIC) or port.
type NetworkInterface struct {
	ID         string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	AssetID    string `gorm:"index" json:"asset_id"`
	NICName    string `gorm:"index" json:"nic_name"`    // Grouping field, e.g. "NIC 1" or "Motherboard NIC"
	Name       string `gorm:"index" json:"name"`        // Port name, e.g. "NIC 1 - Port 1"
	Type       string `json:"type"`                     // "ethernet", "management", "console", "sfp+", "fiber"
	MACAddress  string `json:"mac_address"`                             // e.g. "00:1A:2B:3C:4D:5E"
	IPv4Address string `gorm:"column:ipv4_address" json:"ipv4_address"` // e.g. "192.168.10.15"
	IPv6Address string `gorm:"column:ipv6_address" json:"ipv6_address"` // e.g. "2001:db8::1"
	Speed       string `json:"speed"`                                   // e.g. "1 Gbps", "10 Gbps", "9600 bps"
	MTU        int    `gorm:"default:1500" json:"mtu"`  // e.g. 1500 (Standard) or 9000 (Jumbo Frames)
	VLAN       string `json:"vlan"`                     // e.g. "VLAN 100"
	Status     string `gorm:"default:up" json:"status"` // "up", "down"
}

func (m *NetworkInterface) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// AuditLog represents a record of modifications made to CMDB assets.
type AuditLog struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	AssetID   string    `gorm:"index" json:"asset_id"`
	Action    string    `json:"action"`  // "create", "update", "delete"
	Changes   string    `json:"changes"` // JSON string or summary of modifications
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

func (m *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// AssetRelationship represents uplink or downlink connections between assets (Phase 8/Plan)
type AssetRelationship struct {
	ID       string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	SourceID string `gorm:"index" json:"source_id"` // Asset initiating the link (e.g. child switch)
	TargetID string `gorm:"index" json:"target_id"` // Asset targeted by the link (e.g. parent router)
	Type     string `json:"type"`                   // "uplink" or "downlink"
}

func (m *AssetRelationship) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Manufacturer represents a standardized physical brand/vendor. (Phase 2 Catalog Update)
type Manufacturer struct {
	ID      string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name    string `gorm:"uniqueIndex" json:"name"`
	Website string `json:"website"`
}

func (m *Manufacturer) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// SubGroup represents a dynamic sub-group card portal under a parent Category (Phase 1 Admin Categories)
type SubGroup struct {
	ID          string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CategoryID  string `gorm:"index" json:"category_id"`
	Name        string `json:"name"`                                 // e.g., "Router", "Switch (L2)"
	Icon        string `json:"icon"`                                 // e.g., "i-heroicons-cpu-chip"
	Description string `json:"description"`                          // e.g., "BGP edge routers..."
	Order       int    `gorm:"column:item_order;default:0" json:"order"` // Sorting Weight (First Priority, Phase 1 Sorting)
}

func (m *SubGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Category represents a multi-select category tagging taxonomy. (Phase 2 Catalog Update)
type Category struct {
	ID               string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name             string     `gorm:"uniqueIndex" json:"name"`
	Icon             string     `json:"icon"`                                 // e.g., "i-heroicons-server"
	Description      string     `json:"description"`                          // e.g., "Compute server racks..."
	Order            int        `gorm:"column:item_order;default:0" json:"order"` // Sorting Weight (First Priority, Phase 1 Sorting)
	ParentDependency string     `json:"parent_dependency"`                     // Prerequisite Parent Category (Phase 1 Taxonomy Chains)
	SubGroups        []SubGroup `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"sub_groups,omitempty"`
}

func (m *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// DeviceModel represents a standardized catalog model of hardware devices (Plan & Phase 2 Catalog Update)
type DeviceModel struct {
	ID             string       `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ManufacturerID string       `gorm:"index" json:"manufacturer_id"`
	Manufacturer   Manufacturer `gorm:"foreignKey:ManufacturerID" json:"manufacturer"`
	ModelName      string       `gorm:"index" json:"model_name"`
	Categories     []Category   `gorm:"many2many:device_model_categories;" json:"categories"`
	GeneralInfo    string       `json:"general_info"`
	Subtype        string       `json:"subtype"` // e.g. "Switch (L3)", "Router", "Server", etc.
	IsImported     bool         `gorm:"column:is_imported;default:false" json:"is_imported"`
	Revision       int          `gorm:"default:1" json:"revision"`
	Ports          JSONMap      `gorm:"type:text" json:"ports"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"created_at"`
}

func (m *DeviceModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Revision == 0 {
		m.Revision = 1
	}
	return
}

// DeviceModelRevision represents a historical version of a DeviceModel.
type DeviceModelRevision struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DeviceModelID string    `gorm:"index" json:"device_model_id"`
	Revision      int       `json:"revision"`
	ModelName     string    `json:"model_name"`
	GeneralInfo   string    `json:"general_info"`
	Subtype       string    `json:"subtype"`
	IsImported    bool      `gorm:"column:is_imported" json:"is_imported"`
	Ports         JSONMap   `gorm:"type:text" json:"ports"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (m *DeviceModelRevision) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Permission defines a fine-grained capability tag. (Phase 1 Advanced RBAC)
type Permission struct {
	ID          string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string `gorm:"uniqueIndex" json:"name"` // e.g. "asset:write", "asset:delete"
	Effect      string `json:"effect"`                  // "allow" or "deny"
	Description string `json:"description"`
}

func (m *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Group represents an Access Group. (Phase 1 Advanced RBAC)
type Group struct {
	ID          string       `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string       `gorm:"uniqueIndex" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:group_permissions;" json:"permissions"` // Group-level policies
}

func (m *Group) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// User represents a system operator. (Phase 1 Advanced RBAC)
type User struct {
	ID          string       `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username    string       `gorm:"uniqueIndex" json:"username"`
	Password    string       `json:"-"`
	Role        string       `json:"role"` // For backward compatibility: admin, operator, viewer
	GroupID     *string      `gorm:"index" json:"group_id"`
	Group       *Group       `gorm:"foreignKey:GroupID" json:"group"`
	Permissions []Permission `gorm:"many2many:user_permissions;" json:"permissions"` // User-level overrides!
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// CustomField defines user-defined metadata properties. (Phase 4 Custom Fields)
type CustomField struct {
	ID              string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name            string `gorm:"uniqueIndex" json:"name"` // e.g. "rack_row", "warranty_expiry"
	FieldType       string `json:"field_type"`              // "text", "number", "boolean"
	AssetType       string `gorm:"index" json:"asset_type"` // Server, Switch, Database, Application, etc.
	IsRequired      bool   `gorm:"default:false" json:"is_required"`
	TabGroup        string `json:"tab_group"` // Layout Tab Group Folder (Phase 1 Custom Field Tabs)
	ValidationRegex string `json:"validation_regex"` // Regex Pattern Validation (Phase 1 Custom Field Regex)
}

func (m *CustomField) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// Webhook defines outbound event broadcast webhook registrations. (Phase 3 Webhooks)
type Webhook struct {
	ID    string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	URL   string `gorm:"index" json:"url"`   // target callback URL
	Event string `gorm:"index" json:"event"` // "asset:create", "asset:update", "asset:delete", etc.
}

func (m *Webhook) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}

// SchemaMigration represents a database migration tracking ledger. (Phase 1 Migrations)
type SchemaMigration struct {
	Version   int       `gorm:"primaryKey" json:"version"`
	AppliedAt time.Time `gorm:"autoCreateTime" json:"applied_at"`
}

// ComputeMaintenanceStatus calculates the dynamic maintenance state based on associated windows.
// Priority: in-progress > overdue > coming > none.
func (a *Asset) ComputeMaintenanceStatus() string {
	now := time.Now().UTC()
	hasOverdue := false
	hasComing := false

	for _, w := range a.MaintenanceWindows {
		wStart := w.StartTime.UTC()
		wEnd := w.EndTime.UTC()

		// If the window is currently active and status is "scheduled"
		if w.Status == "scheduled" && now.After(wStart) && now.Before(wEnd) {
			return "in-progress"
		}

		// If the window is scheduled but the EndTime has passed, it is "overdue"
		if w.Status == "scheduled" && now.After(wEnd) {
			hasOverdue = true
			// If it starts in the future, it is "coming"
		} else if wStart.After(now) {
			hasComing = true
		}
	}

	if hasOverdue {
		return "overdue"
	}
	if hasComing {
		return "coming"
	}
	return "none"
}

// LicenseAgreement represents an admin user's acceptance of a specific license version.
type LicenseAgreement struct {
	ID         string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID     string `gorm:"index" json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Version    string `json:"version"`
	AcceptedAt int64  `json:"accepted_at"`
}

func (m *LicenseAgreement) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
