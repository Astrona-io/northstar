package handlers

import (
	"cmdb-backend/internal/models"
	"encoding/json"
)

type ServerProperties struct {
	Manufacturer      string `json:"manufacturer,omitempty"`
	Model             string `json:"model,omitempty"`
	SerialNumber      string `json:"serial_number,omitempty"`
	FormFactor        string `json:"form_factor,omitempty"`
	CPUModel          string `json:"cpu_model,omitempty"`
	CPUSockets        any    `json:"cpu_sockets,omitempty"`
	CPUCores          any    `json:"cpu_cores,omitempty"`
	RAMCapacityGB     any    `json:"ram_capacity_gb,omitempty"`
	RAMType           string `json:"ram_type,omitempty"`
	StorageSubsystem  string `json:"storage_subsystem,omitempty"`
	StorageController string `json:"storage_controller,omitempty"`
	PCIECards         string `json:"pcie_cards,omitempty"`
	NetworkInterfaces string `json:"network_interfaces,omitempty"`
	OSName            string `json:"os_name,omitempty"`
	OSVersion         string `json:"os_version,omitempty"`
}

type DatabaseProperties struct {
	Engine             string `json:"engine,omitempty"`
	Version            string `json:"version,omitempty"`
	Edition            string `json:"edition,omitempty"`
	StorageEngine      string `json:"storage_engine,omitempty"`
	HostCluster        string `json:"host_cluster,omitempty"`
	Environment        string `json:"environment,omitempty"`
	AllocatedStorageGB any    `json:"allocated_storage_gb,omitempty"`
	HASetup            string `json:"ha_setup,omitempty"`
}

type NetworkProperties struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	FormFactor   string `json:"form_factor,omitempty"`
	Firmware     string `json:"firmware,omitempty"`
	ManagementIP string `json:"management_ip,omitempty"`
	MACAddress   string `json:"mac_address,omitempty"`
	NetworkRole  string `json:"network_role,omitempty"`
	PortConfig   string `json:"port_config,omitempty"`
	Throughput   string `json:"throughput,omitempty"`
}

type ApplicationProperties struct {
	Vendor       string `json:"vendor,omitempty"`
	Version      string `json:"version,omitempty"`
	Environment  string `json:"environment,omitempty"`
	Criticality  string `json:"criticality,omitempty"`
	Framework    string `json:"framework,omitempty"`
	FrontendTech string `json:"frontend_tech,omitempty"`
	URL          string `json:"url,omitempty"`
	AuthMethod   string `json:"auth_method,omitempty"`
}

func sanitizeProperties(assetType string, raw models.JSONMap) (models.JSONMap, error) {
	if raw == nil {
		return models.JSONMap{}, nil
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	var target any
	switch assetType {
	case "Server":
		target = &ServerProperties{}
	case "Database":
		target = &DatabaseProperties{}
	case "Router", "Switch":
		target = &NetworkProperties{}
	case "Application":
		target = &ApplicationProperties{}
	default:
		// If type is unknown, just return the raw map, or we could return empty.
		// Returning raw allows custom types to have custom properties.
		return cleanEmptyStrings(raw), nil
	}

	if err := json.Unmarshal(b, target); err != nil {
		return nil, err
	}

	b2, _ := json.Marshal(target)
	var cleanMap models.JSONMap
	json.Unmarshal(b2, &cleanMap)

	return cleanEmptyStrings(cleanMap), nil
}

func cleanEmptyStrings(m models.JSONMap) models.JSONMap {
	clean := make(models.JSONMap)
	for k, v := range m {
		if s, ok := v.(string); ok && s == "" {
			continue // skip empty strings
		}
		clean[k] = v
	}
	return clean
}
