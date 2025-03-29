package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Availability struct {
	DayOfWeek  int    `json:"dayOfWeek"`  // 0-6 (Sunday-Saturday)
	StartTime  string `json:"startTime"`   // Format: "HH:MM"
	EndTime    string `json:"endTime"`     // Format: "HH:MM"
	IsAvailable bool  `json:"isAvailable"`
}

type MentorDetails struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null"`
	User          User          `gorm:"foreignKey:UserID"`
	Role          string        `gorm:"type:varchar(20);not null;default:'mentor'"`
	Experience    string        `gorm:"type:text"`
	Skills        []string      `gorm:"type:text[]"`  // Primary skills array
	Certifications []string      `gorm:"type:text[]"`
	Availability  []Availability `gorm:"-"`                              // Stored as JSON in AvailabilityJSON
	AvailabilityJSON string     `gorm:"type:jsonb;column:availability"` // Internal storage field
	Rating         float64      `gorm:"default:0"`
	ReviewsCount   int          `gorm:"default:0"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Tags for detailed categorization (skills, expertise, interests, etc.)
	Tags []Tag `gorm:"many2many:mentor_tags;"`
}

func (m *MentorDetails) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.Role == "" {
		m.Role = "mentor"
	}
	return nil
}

// BeforeSave handles JSON conversion for availability
func (m *MentorDetails) BeforeSave(tx *gorm.DB) error {
	if len(m.Availability) > 0 {
		data, err := json.Marshal(m.Availability)
		if err != nil {
			return err
		}
		m.AvailabilityJSON = string(data)
	}
	return nil
}

// AfterFind handles JSON parsing for availability
func (m *MentorDetails) AfterFind(tx *gorm.DB) error {
	if m.AvailabilityJSON != "" {
		return json.Unmarshal([]byte(m.AvailabilityJSON), &m.Availability)
	}
	return nil
}
