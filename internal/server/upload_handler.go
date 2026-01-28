package server

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"sentinel/internal/database"
	"sentinel/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) UploadHandler(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int(val.(uint))

	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file received"})
		return
	}

	filename := filepath.Base(file.Filename)
	dst := "./uploads/" + filename

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	urls, err := processFile(dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error processing file: %s", err)})
		return
	}

	// Concurrent
	go func(urlList []string, uid int, fPath string) {
		for _, u := range urlList {
			cleanU := strings.TrimSpace(u)
			if isValidURL(cleanU) {
				job := models.Job{
					URL:      cleanU,
					UserID:   uid,
					Status:   "pending",
					FilePath: fPath,
					JobType:  "web",
				}
				if err := database.CreateJob(s.WorkerPool.DB, &job); err == nil {
					s.WorkerPool.JobChan <- job
				}
			}
		}
	}(urls, userID, dst)

	c.JSON(http.StatusAccepted, gin.H{
		"message":     "File uploaded successfully. Processing in background.",
		"filename":    filename,
		"total_found": len(urls),
	})
}

func isValidURL(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
