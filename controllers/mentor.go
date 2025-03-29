package controllers

import (
	"mentorship-backend/config"
	"mentorship-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MentorController struct{}

func NewMentorController() *MentorController {
	return &MentorController{}
}

// CreateMentorProfile creates or updates mentor profile
func (mc *MentorController) CreateMentorProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var mentorDetails models.MentorDetails
	if err := c.ShouldBindJSON(&mentorDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the UserID from the authenticated user
	mentorDetails.UserID = userID.(uuid.UUID)
	mentorDetails.Role = "mentor" // Explicitly set role

	// Check if mentor profile already exists
	var existingProfile models.MentorDetails
	result := config.GetDB().Where("user_id = ?", userID).First(&existingProfile)
	
	if result.Error == nil {
		// Update existing profile
		existingProfile.Experience = mentorDetails.Experience
		existingProfile.Skills = mentorDetails.Skills
		existingProfile.Certifications = mentorDetails.Certifications
		existingProfile.Availability = mentorDetails.Availability
		
		if err := config.GetDB().Save(&existingProfile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update mentor profile"})
			return
		}
		c.JSON(http.StatusOK, existingProfile)
		return
	}

	// Create new profile
	if err := config.GetDB().Create(&mentorDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mentor profile"})
		return
	}

	c.JSON(http.StatusCreated, mentorDetails)
}

// GetMentorProfile gets mentor profile by ID
func (mc *MentorController) GetMentorProfile(c *gin.Context) {
	mentorID := c.Param("id")
	
	var mentorDetails models.MentorDetails
	if err := config.GetDB().Preload("User").Preload("Tags").First(&mentorDetails, "id = ?", mentorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mentor profile not found"})
		return
	}

	c.JSON(http.StatusOK, mentorDetails)
}

// ListMentors lists all mentors with optional filters
func (mc *MentorController) ListMentors(c *gin.Context) {
	var mentors []models.MentorDetails
	
	query := config.GetDB().Preload("User").Preload("Tags")
	
	// Add skill filter if provided
	if skill := c.Query("skill"); skill != "" {
		query = query.Where("? = ANY(skills)", skill)
	}

	// Add tag filter if provided
	if tagName := c.Query("tag"); tagName != "" {
		query = query.Joins("JOIN mentor_tags ON mentor_tags.mentor_details_id = mentor_details.id").
			Joins("JOIN tags ON tags.id = mentor_tags.tag_id").
			Where("tags.name = ?", tagName)
	}

	// Execute query
	if err := query.Find(&mentors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch mentors"})
		return
	}

	c.JSON(http.StatusOK, mentors)
}

// UpdateAvailability updates mentor's availability
func (mc *MentorController) UpdateAvailability(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var availability []models.Availability
	if err := c.ShouldBindJSON(&availability); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mentorDetails models.MentorDetails
	if err := config.GetDB().Where("user_id = ?", userID).First(&mentorDetails).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mentor profile not found"})
		return
	}

	mentorDetails.Availability = availability
	if err := config.GetDB().Save(&mentorDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Availability updated successfully"})
}
