package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Follow struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null"`  // User who is following
	Follower    User      `gorm:"foreignKey:FollowerID"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null"`  // Mentor being followed
	Following   User      `gorm:"foreignKey:FollowingID"`
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (f *Follow) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
