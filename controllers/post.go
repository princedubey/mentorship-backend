package controllers

import (
	"mentorship-backend/config"
	"mentorship-backend/models"
	"mentorship-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// CreatePost creates a new post
func (pc *PostController) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload if present
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		// Upload to Cloudinary
		url, err := utils.UploadImage(file, "posts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		
		// Add the URL to the post's media URLs
		post.MediaURLs = append(post.MediaURLs, url)
	}

	post.UserID = userID.(uuid.UUID)

	if err := config.GetDB().Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPost gets a post by ID
func (pc *PostController) GetPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Post ID is required"})
		return
	}

	var post models.Post
	if err := config.GetDB().Preload("User").Preload("Comments").Preload("Tags").Preload("SavedBy").Preload("Likes").First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	// Get likes and comments count
	likesCount, _ := post.GetLikesCount(config.GetDB())
	commentsCount, _ := post.GetCommentsCount(config.GetDB())

	// Update analytics
	post.Analytics.Likes = int(likesCount)
	post.Analytics.CommentCount = int(commentsCount)

	c.JSON(200, post)
}

// ListPosts lists all posts with optional filters and search
func (pc *PostController) ListPosts(c *gin.Context) {
	var posts []models.Post
	query := config.GetDB().Preload("User").
		Preload("Tags").
		Preload("Comments", "parent_id IS NULL").
		Preload("OriginalPost").
		Preload("OriginalPost.User")

	// Add tag filter
	if tagName := c.Query("tag"); tagName != "" {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", tagName)
	}

	// Add user filter
	if userID := c.Query("user"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Add search filter
	if search := c.Query("search"); search != "" {
		searchTerms := strings.Split(search, " ")
		for _, term := range searchTerms {
			query = query.Where("content ILIKE ?", "%"+term+"%")
		}
	}

	// Add date range filter
	if startDate := c.Query("startDate"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
	}
	if endDate := c.Query("endDate"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
	}

	// Only show public posts for non-owners
	currentUserID, exists := c.Get("userID")
	if !exists {
		query = query.Where("is_private = ?", false)
	} else {
		query = query.Where("is_private = ? OR user_id = ?", false, currentUserID)
	}

	if err := query.Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// SharePost shares an existing post
func (pc *PostController) SharePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	originalPostID := c.Param("id")
	var originalPost models.Post
	if err := config.GetDB().First(&originalPost, "id = ?", originalPostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Original post not found"})
		return
	}

	// Create shared post
	sharedPost := models.Post{
		UserID:         userID.(uuid.UUID),
		OriginalPostID: &originalPost.ID,
		IsPrivate:      false,
	}

	tx := config.GetDB().Begin()
	if err := tx.Create(&sharedPost).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share post"})
		return
	}

	// Increment share count
	if err := tx.Model(&originalPost).
		UpdateColumn("analytics__shares", gorm.Expr("analytics__shares + ?", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update share count"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, sharedPost)
}

// SavePost allows a user to save/bookmark a post
func (pc *PostController) SavePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := config.GetDB().First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	tx := config.GetDB().Begin()
	// Add user to SavedBy
	if err := tx.Model(&post).Association("SavedBy").Append(&models.User{ID: userID.(uuid.UUID)}); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save post"})
		return
	}

	// Increment saved count
	if err := tx.Model(&post).
		UpdateColumn("analytics__saved_count", gorm.Expr("analytics__saved_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update saved count"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Post saved successfully"})
}

// GetPostAnalytics gets analytics for a post
func (pc *PostController) GetPostAnalytics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := config.GetDB().First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Only post owner can see analytics
	if post.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view analytics"})
		return
	}

	likesCount, _ := post.GetLikesCount(config.GetDB())
	commentsCount, _ := post.GetCommentsCount(config.GetDB())

	// Update analytics
	post.Analytics.Likes = int(likesCount)
	post.Analytics.CommentCount = int(commentsCount)

	analytics := struct {
		models.PostAnalytics
		EngagementRate float64 `json:"engagementRate"`
	}{
		PostAnalytics: post.Analytics,
		EngagementRate: float64(likesCount + int64(post.Analytics.CommentCount) +
			int64(post.Analytics.Shares) + int64(post.Analytics.SavedCount)) / float64(post.Analytics.Views) * 100,
	}

	c.JSON(http.StatusOK, analytics)
}

// AddTagsToPost adds tags to a post
func (pc *PostController) AddTagsToPost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := config.GetDB().First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Verify post ownership
	if post.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to modify this post"})
		return
	}

	var tagIDs []uuid.UUID
	if err := c.ShouldBindJSON(&tagIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tags []models.Tag
	if err := config.GetDB().Find(&tags, "id IN ?", tagIDs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag IDs"})
		return
	}

	if err := config.GetDB().Model(&post).Association("Tags").Append(tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags added successfully"})
}

// DeletePost deletes a post and its associated images
func (pc *PostController) DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := config.GetDB().First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user is authorized to delete the post
	if post.UserID != userID.(uuid.UUID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this post"})
		return
	}

	tx := config.GetDB().Begin()

	// Delete associated likes
	if err := tx.Where("post_id = ?", postID).Delete(&models.Like{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post likes"})
		return
	}

	// Delete associated comments
	if err := tx.Where("post_id = ?", postID).Delete(&models.Comment{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post comments"})
		return
	}

	// Delete associated tags
	if err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", postID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post tags"})
		return
	}

	// Delete associated saves
	if err := tx.Exec("DELETE FROM user_saved_posts WHERE post_id = ?", postID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post saves"})
		return
	}

	// Delete the post
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	tx.Commit()

	// Delete images from Cloudinary in a separate goroutine
	go func() {
		if err := utils.DeleteImagesFromPost(post.MediaURLs); err != nil {
			// Log the error but don't fail the deletion
			log.Printf("Failed to delete images from Cloudinary: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
