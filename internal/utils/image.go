package utils

import (
	"path/filepath"
	"strings"
)

// SupportedImageExt returns true if the file has a supported image extension
func SupportedImageExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp":
		return true
	default:
		return false
	}
}
