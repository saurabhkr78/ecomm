package domain

import "time"

type Cart struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint      `json:"user_id"`
	ProductId uint      `json:"product_id"`
	SellerId  uint      `json:"seller_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	ImageUrl  string    `json:"image_url"`
	Qty       int       `json:"qty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
