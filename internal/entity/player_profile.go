package entity

import "time"

type PlayerProfile struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint64     `gorm:"not null" json:"userId"`
	Level             string     `gorm:"type:enum('beginner','average','good','pro');default:'average'" json:"level"`
	Gender            *string    `gorm:"type:enum('male','female','other')" json:"gender,omitempty"`
	PreferredLocation *string    `gorm:"size:255" json:"preferredLocation,omitempty"`
	Latitude          *float64   `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
	Longitude         *float64   `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
	Bio               *string    `gorm:"type:text" json:"bio,omitempty"`
	Distance          *float64   `gorm:"-" json:"distance,omitempty"` // Calculated field, not stored in DB
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt         *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
