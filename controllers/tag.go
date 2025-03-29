package controllers

import (
	"mentorship-backend/config"
	"mentorship-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagController struct{}

func NewTagController() *TagController {
	return &TagController{}
}

// CreateTag creates a new tag
func (tc *TagController) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.GetDB().Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// ListTags lists all tags with optional category filter
func (tc *TagController) ListTags(c *gin.Context) {
	var tags []models.Tag

	query := config.GetDB()
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// AddTagsToUser adds tags to a user
func (tc *TagController) AddTagsToUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var tagIDs []uuid.UUID
	if err := c.ShouldBindJSON(&tagIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.GetDB().First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var tags []models.Tag
	if err := config.GetDB().Find(&tags, "id IN ?", tagIDs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag IDs"})
		return
	}

	if err := config.GetDB().Model(&user).Association("Tags").Append(tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags added successfully"})
}

// AddTagsToMentor adds tags to a mentor
func (tc *TagController) AddTagsToMentor(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var tagIDs []uuid.UUID
	if err := c.ShouldBindJSON(&tagIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var mentor models.MentorDetails
	if err := config.GetDB().First(&mentor, "user_id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mentor profile not found"})
		return
	}

	var tags []models.Tag
	if err := config.GetDB().Find(&tags, "id IN ?", tagIDs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag IDs"})
		return
	}

	if err := config.GetDB().Model(&mentor).Association("Tags").Append(tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags added successfully"})
}
