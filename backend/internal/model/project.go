package model

import "time"

// RenovationProject 装修项目。
type RenovationProject struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Name             string     `gorm:"size:120;not null" json:"name"`
	HouseType        string     `gorm:"size:32;not null" json:"house_type"`
	Area             float64    `gorm:"type:decimal(10,2);not null" json:"area"`
	DecorStyle       string     `gorm:"size:32;not null" json:"decor_style"`
	Address          string     `gorm:"size:255" json:"address"`
	OwnerID          uint       `gorm:"not null;index" json:"owner_id"`
	DesignerID       uint       `gorm:"index" json:"designer_id"`
	ForemanID        uint       `gorm:"index" json:"foreman_id"`
	Status           string     `gorm:"size:32;not null;index" json:"status"`
	ContractAmount   float64    `gorm:"type:decimal(14,2);not null;default:0" json:"contract_amount"`
	StartDate        *time.Time `gorm:"type:date" json:"start_date"`
	ExpectedEndDate  *time.Time `gorm:"type:date" json:"expected_end_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (RenovationProject) TableName() string { return "renovation_projects" }
