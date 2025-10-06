package domain

type BankAccount struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id"`
	BankAccount uint   `json:"bank_account" gorm:"index;unique;not null"`
	IFSCCode    string `json:"ifsc_code"`
	PaymentType string `json:"payment_type"`
	CreatedAt   string `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   string `json:"updated_at" gorm:"default:current_timestamp"`
}
