package common

import (
	"database/sql"
	"log"
)

// Role constants
const (
	RoleAdmin            = "admin"
	RoleManager          = "manager"
	RolePlacementOfficer = "placement_officer"
)

// ValidRoles defines the allowed roles in the system
var ValidRoles = []string{RoleAdmin, RoleManager, RolePlacementOfficer}

// SystemSettings holds system configuration
type SystemSettings struct {
	ApprovalMode string // "auto" or "manual"
}

// Default settings
var DefaultSettings = SystemSettings{
	ApprovalMode: "manual", // Default to manual approval
}

// Global settings instance
var CurrentSettings SystemSettings

// InitSystemSettings initializes the system settings from the database
func InitSystemSettings(db *sql.DB) {
	// Try to load settings from database
	var approvalMode string
	err := db.QueryRow("SELECT approval_mode FROM system_settings LIMIT 1").Scan(&approvalMode)
	if err != nil {
		if err == sql.ErrNoRows {
			// No settings found, use defaults
			log.Println("No system settings found in database, using defaults")
			CurrentSettings = DefaultSettings

			// Optionally insert default settings
			_, err = db.Exec("INSERT INTO system_settings (approval_mode) VALUES (?)", DefaultSettings.ApprovalMode)
			if err != nil {
				log.Printf("Error inserting default settings: %v", err)
			}
		} else {
			// Error occurred, use defaults
			log.Printf("Error loading system settings: %v", err)
			CurrentSettings = DefaultSettings
		}
	} else {
		// Settings found, use them
		CurrentSettings = SystemSettings{
			ApprovalMode: approvalMode,
		}
		log.Printf("Loaded system settings: approval mode = %s", approvalMode)
	}
}
