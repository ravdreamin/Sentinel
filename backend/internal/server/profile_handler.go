package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) ProfileHandler(c *gin.Context) {
	
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your protected profile!",
		"user_id": userID,
	})
}
