package models

import (
	"time"

	"gorm.io/gorm"
)

// Customer struct represents a customer in the system
type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`             // Mandatory name
	Email     string         `gorm:"unique;default:null" json:"email"` // Optional email (nullable)
	Phone     string         `gorm:"not null" json:"phone"`            // Mandatory phone number for WhatsApp
	Address   string         `gorm:"default:null" json:"address"`      // Optional address (nullable)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
