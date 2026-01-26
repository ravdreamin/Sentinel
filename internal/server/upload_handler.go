package server

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
	})

}
