package entity

import "time"

type Notification struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null" json:"userId"`
	Title     *string    `gorm:"size:255" json:"title,omitempty"`
	Message   *string    `gorm:"type:text" json:"message,omitempty"`
	IsRead    bool       `gorm:"default:false" json:"isRead"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
