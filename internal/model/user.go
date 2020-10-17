package model

import "time"

// User represents the user model
// swagger:model
type User struct {
	Base
	FirstName string `json:"first_name" gorm:"type:varchar(255)"`
	LastName  string `json:"last_name" gorm:"type:varchar(255)"`
	Email     string `json:"email" gorm:"type:varchar(255)"`
	Mobile    string `json:"mobile,omitempty" gorm:"type:varchar(255)"`

	Username     string     `json:"username" gorm:"type:varchar(255);unique_index;not null"`
	Password     string     `json:"-" gorm:"type:varchar(255);not null"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	Blocked      bool       `json:"blocked" gorm:"not null;default:false"`
	RefreshToken string     `json:"-" gorm:"type:varchar(255);unique_index"`

	Role string `json:"role" gorm:"varchar(255)"`
}
