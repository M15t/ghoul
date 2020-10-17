package model

// Country represents the country model
// swagger:model
type Country struct {
	Base
	Name      string `json:"name" gorm:"type:varchar(255)"`
	Code      string `json:"code" gorm:"type:varchar(10)"`
	PhoneCode string `json:"phone_code" gorm:"type:varchar(10)"`
}
