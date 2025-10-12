package entity

import "time"

type Feedback struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID  uint64     `gorm:"not null" json:"bookingId"`
	FromUserID uint64     `gorm:"not null" json:"fromUserId"`
	ToUserID   uint64     `gorm:"not null" json:"toUserId"`
	Rating     *int       `gorm:"check:rating BETWEEN 1 AND 5" json:"rating,omitempty"`
	Comment    *string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt  *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Booking  *Booking `gorm:"foreignKey:BookingID" json:"booking,omitempty"`
	FromUser *User    `gorm:"foreignKey:FromUserID" json:"fromUser,omitempty"`
	ToUser   *User    `gorm:"foreignKey:ToUserID" json:"toUser,omitempty"`
}
