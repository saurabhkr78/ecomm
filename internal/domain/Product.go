package domain

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint           `json:"id" gorm:"PrimaryKey"`
	Name        string         `json:"name" gorm:"index;not null"`
	Description string         `json:"description"`
	CategoryId  string         `json:"category_id"`
	ImageUrl    string         `json:"image_url"`
	Price       float64        `json:"price"`
	UserId      int            `json:"user_id"`
	Stock       int            `json:"stock"`
	CreatedAt   time.Time      `json:"created_at" default:"current_timestamp"`
	UpdatedAt   time.Time      `json:"updated_at" default:"current_timestamp"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
