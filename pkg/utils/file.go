package utils

import (
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func IsAllowedImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	}
	return false
}

func GenerateUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	uniqueID := uuid.New().String()
	return uniqueID + ext
}