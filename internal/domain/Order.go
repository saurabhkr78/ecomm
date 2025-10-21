package domain

import "time"

type Order struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	UserID         uint        `json:"user_id"`
	Status         string      `json:"status"`
	TotalAmount    float64     `json:"total_amount"`
	TransactionId  string      `json:"transaction_id"`
	OrderRefNumber int         `json:"order_ref_number"`
	PaymentId      string      `json:"payment_id"`
	Items          []OrderItem `json:"items"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
