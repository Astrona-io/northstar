package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type LatencyPingRes struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	IPAddress  string  `json:"ip_address"`
	Status     string  `json:"status"` // "online", "offline", "maintenance"
	LatencyMS  int     `json:"latency_ms"`
	PacketLoss float64 `json:"packet_loss"`
}

// GetLatencyPing handles GET /api/monitoring/ping (Phase 1 Latency Monitoring)
func GetLatencyPing(c echo.Context) error {
	var assets []models.Asset
	if err := database.DB.Find(&assets).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to query assets: " + err.Error()})
	}

	rand.Seed(time.Now().UnixNano())

	var response []LatencyPingRes
	for _, asset := range assets {
		ip := ""
		if asset.IPAddress != nil {
			ip = *asset.IPAddress
		}

		status := "online"
		latency := rand.Intn(26) + 4 // 4ms to 30ms randomized
		loss := 0.0

		if asset.Status == "inactive" {
			status = "offline"
			latency = 0
			loss = 100.0
		} else if asset.Status == "maintenance" {
			status = "maintenance"
			latency = rand.Intn(45) + 35 // 35ms to 80ms under load
			loss = rand.Float64() * 5.0  // up to 5% loss
		}

		response = append(response, LatencyPingRes{
			ID:         asset.ID,
			Name:       asset.Name,
			IPAddress:  ip,
			Status:     status,
			LatencyMS:  latency,
			PacketLoss: loss,
		})
	}

	return c.JSON(http.StatusOK, response)
}
