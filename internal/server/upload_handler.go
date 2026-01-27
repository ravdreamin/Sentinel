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
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no file recieve",
		})
		return
	}

	filename := filepath.Base(file.Filename)
	dst := "./uploads/" + filename

	if err := c.SaveUploadedFile(file, dst); err != nil {
		fmt.Println("Error saving file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	urls, err := processFile(dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Error processing file: %s", err),
		})
		return
	}

	validURLs := []string{}
	for _, u := range urls {
		cleanU := strings.TrimSpace(u)

		if isValidURL(cleanU) {
			validURLs = append(validURLs, cleanU)
			job := models.Job{
				URL: cleanU,
			}

			err := database.CreateJob(s.WorkerPool.DB, &job)
			if err != nil {
				fmt.Println("Failed to save job:", err)
				continue
			}

			s.WorkerPool.JobChan <- job

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "File uploaded successfully",
		"filename":    filename,
		"total_found": len(urls),
		"valid_urls":  len(validURLs),
		"urls":        urls,
	})

}

func isValidURL(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
