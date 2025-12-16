package handler

import (
	"go-ecommerce/pkg/utils"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if !utils.IsAllowedImage(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed. Only jpg, jpeg, png, webp"})
		return
	}

	newFilename := utils.GenerateUniqueFilename(file.Filename)
	savePath := filepath.Join("uploads", newFilename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileURL := "http://localhost:8080/uploads/" + newFilename
	
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fileURL,
	})
}