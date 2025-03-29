package controllers

import (
	"mentorship-backend/config"
	"mentorship-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

// CreateComment creates a new comment
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.UserID = userID.(uuid.UUID)
	postUUID, err := uuid.Parse(postID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	comment.PostID = postUUID

	tx := config.GetDB().Begin()
	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Increment comment count
	if err := tx.Model(&models.Post{}).Where("id = ?", postID).
		UpdateColumn("analytics__comment_count", gorm.Expr("analytics__comment_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment count"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, comment)
}

// GetComments gets all comments for a post
func (cc *CommentController) GetComments(c *gin.Context) {
	postID := c.Param("id")
	
	var comments []models.Comment
	if err := config.GetDB().Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("User").
		Preload("Replies.User").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// ReplyToComment creates a reply to a comment
func (cc *CommentController) ReplyToComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commentID := c.Param("id")
	var reply models.Comment
	if err := c.ShouldBindJSON(&reply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get parent comment to get post ID
	var parentComment models.Comment
	if err := config.GetDB().First(&parentComment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
		return
	}

	commentUUID, _ := uuid.Parse(commentID)
	reply.UserID = userID.(uuid.UUID)
	reply.PostID = parentComment.PostID
	reply.ParentID = &commentUUID

	tx := config.GetDB().Begin()
	if err := tx.Create(&reply).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
		return
	}

	// Increment comment count
	if err := tx.Model(&models.Post{}).Where("id = ?", parentComment.PostID).
		UpdateColumn("analytics__comment_count", gorm.Expr("analytics__comment_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment count"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, reply)
}
