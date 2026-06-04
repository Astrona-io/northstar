package models

import (
	"testing"
	"time"
)

func TestComputeMaintenanceStatus(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		windows  []MaintenanceWindow
		expected string
	}{
		{
			name:     "No windows",
			windows:  []MaintenanceWindow{},
			expected: "none",
		},
		{
			name: "Completed window in-progress ignored",
			windows: []MaintenanceWindow{
				{
					Status:    "completed",
					StartTime: now.Add(-1 * time.Hour),
					EndTime:   now.Add(1 * time.Hour),
				},
			},
			expected: "none",
		},
		{
			name: "Cancelled window in-progress ignored",
			windows: []MaintenanceWindow{
				{
					Status:    "cancelled",
					StartTime: now.Add(-1 * time.Hour),
					EndTime:   now.Add(1 * time.Hour),
				},
			},
			expected: "none",
		},
		{
			name: "Single scheduled window coming",
			windows: []MaintenanceWindow{
				{
					Status:    "scheduled",
					StartTime: now.Add(1 * time.Hour),
					EndTime:   now.Add(2 * time.Hour),
				},
			},
			expected: "coming",
		},
		{
			name: "Single scheduled window overdue",
			windows: []MaintenanceWindow{
				{
					Status:    "scheduled",
					StartTime: now.Add(-2 * time.Hour),
					EndTime:   now.Add(-1 * time.Hour),
				},
			},
			expected: "overdue",
		},
		{
			name: "Single scheduled window in-progress",
			windows: []MaintenanceWindow{
				{
					Status:    "scheduled",
					StartTime: now.Add(-1 * time.Hour),
					EndTime:   now.Add(1 * time.Hour),
				},
			},
			expected: "in-progress",
		},
		{
			name: "Priority check: in-progress takes precedence over coming and overdue",
			windows: []MaintenanceWindow{
				{
					Status:    "scheduled",
					StartTime: now.Add(-2 * time.Hour),
					EndTime:   now.Add(-1 * time.Hour), // overdue
				},
				{
					Status:    "scheduled",
					StartTime: now.Add(1 * time.Hour),
					EndTime:   now.Add(2 * time.Hour), // coming
				},
				{
					Status:    "scheduled",
					StartTime: now.Add(-30 * time.Minute),
					EndTime:   now.Add(30 * time.Minute), // in-progress
				},
			},
			expected: "in-progress",
		},
		{
			name: "Priority check: overdue takes precedence over coming",
			windows: []MaintenanceWindow{
				{
					Status:    "scheduled",
					StartTime: now.Add(-2 * time.Hour),
					EndTime:   now.Add(-1 * time.Hour), // overdue
				},
				{
					Status:    "scheduled",
					StartTime: now.Add(1 * time.Hour),
					EndTime:   now.Add(2 * time.Hour), // coming
				},
			},
			expected: "overdue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			asset := Asset{
				MaintenanceWindows: tt.windows,
			}
			actual := asset.ComputeMaintenanceStatus()
			if actual != tt.expected {
				t.Errorf("expected maintenance status %q, got %q", tt.expected, actual)
			}
		})
	}
}
