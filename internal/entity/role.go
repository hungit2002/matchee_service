package entity

import "time"

type Role struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	UserRoles []UserRole `gorm:"foreignKey:RoleID" json:"userRoles,omitempty"`
}

const (
	RolePlayer     = "player"
	RoleAdmin      = "admin"
	RoleVenueOwner = "venue_owner"
)
