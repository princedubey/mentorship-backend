package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler returns health status of the application
func HealthCheckHandler(c *gin.Context) {
	status := struct {
		Status      string    `json:"status"`
		Timestamp   time.Time `json:"timestamp"`
		Environment string    `json:"environment"`
		Version     string    `json:"version"`
	}{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Environment: "production",
		Version:     "1.0.0",
	}

	c.JSON(200, status)
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"status":  "error",
		"message": fmt.Sprintf("Endpoint %s not found", c.Request.URL.Path),
	})
}
