package domain

import "time"

type BankAccount struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	BankAccount uint      `json:"bank_account" gorm:"index;unique;not null"`
	IFSCCode    string    `json:"ifsc_code"`
	PaymentType string    `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
