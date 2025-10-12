package entity

import "time"

type MatchGroupMember struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MatchGroupID uint64     `gorm:"not null" json:"matchGroupId"`
	UserID       uint64     `gorm:"not null" json:"userId"`
	IsLeader     bool       `gorm:"default:false" json:"isLeader"`
	JoinedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"joinedAt"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	// Relationships
	MatchGroup *MatchGroup `gorm:"foreignKey:MatchGroupID" json:"matchGroup,omitempty"`
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
