package entity

import "time"

type AvailabilitySlot struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CourtID     uint64     `gorm:"not null" json:"courtId"`
	StartTime   time.Time  `gorm:"not null" json:"startTime"`
	EndTime     time.Time  `gorm:"not null" json:"endTime"`
	IsAvailable bool       `gorm:"default:true" json:"isAvailable"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	Court *Court `gorm:"foreignKey:CourtID" json:"court,omitempty"`
}
