package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sentinel/internal/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) JobStatusHandler(c *gin.Context) {
	filename := c.Param("filename")
	// Use filename as file_path (add ./uploads/ prefix if stored that way)
	// In UploadHandler: dst := "./uploads/" + filename
	filePath := "./uploads/" + filename

	total, completed, failed, err := database.GetJobProgress(s.WorkerPool.DB, filePath)
	if err != nil {
		fmt.Printf("[DEBUG] Progress Err: %v | Path: %s\n", err, filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch status"})
		return
	}

	fmt.Printf("[DEBUG] Status: T=%d C=%d F=%d | Path: %s\n", total, completed, failed, filePath)

	status := "processing"
	if total > 0 && (completed+failed) == total {
		status = "completed"
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"completed": completed,
		"failed":    failed,
		"status":    status,
	})
}

func (s *Server) JobDownloadHandler(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "./uploads/" + filename

	results, err := database.GetJobResults(s.WorkerPool.DB, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch results"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No results found (or job pending)"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s_results.json", filename))
	c.Header("Content-Type", "application/json")

	encoder := json.NewEncoder(c.Writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		// Log error
	}
}

func (s *Server) ListJobsHandler(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int(val.(uint))

	// Guests don't have stored jobs list (per requirement)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{"jobs": []string{}})
		return
	}

	files, err := database.GetUserJobs(s.WorkerPool.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": files})
}

func (s *Server) DeleteJobHandler(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int(val.(uint))

	if userID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Guests cannot delete jobs"})
		return
	}

	filename := c.Param("filename")
	filePath := "./uploads/" + filename

	fmt.Printf("[DEBUG] Delete request: filename=%s filePath=%s userID=%d\n", filename, filePath, userID)

	if err := database.DeleteJobByFilePath(s.WorkerPool.DB, filePath, userID); err != nil {
		fmt.Printf("[DEBUG] Delete error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}

	os.Remove(filePath)
	c.JSON(http.StatusOK, gin.H{"message": "Job deleted"})
}

func (s *Server) JobMetricsHandler(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "./uploads/" + filename

	metrics, err := database.GetJobMetrics(s.WorkerPool.DB, filePath)
	if err != nil {
		fmt.Printf("[DEBUG] Metrics Err: %v | Path: %s\n", err, filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch metrics"})
		return
	}

	fmt.Printf("[DEBUG] Metrics: %+v | Path: %s\n", metrics, filePath)
	c.JSON(http.StatusOK, metrics)
}
