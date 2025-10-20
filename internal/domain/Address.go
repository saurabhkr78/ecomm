package domain

import "time"

type Address struct {
	ID           uint      `gorm:"primaryKey,json:"id"`
	AddressLine1 string    `json:"address_line_1"`
	Street       string    `json:"street"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	PinCode      uint      `json:"pin_code"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
