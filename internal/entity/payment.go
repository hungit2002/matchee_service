package entity

import "time"

type Payment struct {
	ID                  uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID           uint64     `gorm:"not null" json:"bookingId"`
	PayerID             uint64     `gorm:"not null" json:"payerId"`
	Amount              *float64   `gorm:"type:decimal(10,2)" json:"amount,omitempty"`
	Method              *string    `gorm:"type:enum('paypal','momo','zalo','cash','bank_transfer')" json:"method,omitempty"`
	Status              string     `gorm:"type:enum('pending','success','failed');default:'pending'" json:"status"`
	PayPalOrderID       *string    `gorm:"column:paypal_order_id" json:"paypalOrderId,omitempty"`
	PayPalTransactionID *string    `gorm:"column:paypal_transaction_id" json:"paypalTransactionId,omitempty"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt           *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Booking *Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	Payer   *User    `gorm:"foreignKey:PayerID" json:"payer,omitempty"`
}
