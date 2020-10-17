package model

import (
	"time"
)

// Base contains common fields for all models
// Do not use gorm.Model because of uint ID
type Base struct {
	// ID of the record
	ID int `json:"id" gorm:"primary_key"`
	// The time that record is created
	CreatedAt time.Time `json:"created_at"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt *time.Time `json:"deleted_at" sql:"index"`
}
