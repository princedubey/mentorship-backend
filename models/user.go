package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser Role = "user"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirebaseUID  string    `gorm:"type:varchar(128);unique;not null"` // Firebase UID
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex"` // Optional for phone auth
	PhoneNumber  string    `gorm:"type:varchar(20);uniqueIndex"` // Optional for email auth
	Password     string    `gorm:""` // Optional now, as Firebase handles auth
	Role         Role      `gorm:"type:varchar(20);not null;default:'user'"`
	Bio          string    `gorm:"type:text"`
	AvatarURL    string    `gorm:"type:text"`
	IsPrivate    bool      `gorm:"default:false"`
	IsActive     bool      `gorm:"default:true"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relationships
	SavedPosts []Post `gorm:"many2many:user_saved_posts;"`
	Tags       []Tag  `gorm:"many2many:user_tags;"`
}

// BeforeCreate will set default role
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Role == "" {
		u.Role = RoleUser
	}
	return nil
}
