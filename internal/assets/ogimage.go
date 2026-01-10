package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OGImageConfig holds configuration for OG image handling
type OGImageConfig struct {
	ImagePath string // Local path to OG image file
	BaseURL   string // Site base URL for absolute URL generation
}

// ProcessOGImage copies the OG image to output and returns the URL
func ProcessOGImage(config OGImageConfig, outputDir string) (string, error) {
	if config.ImagePath == "" {
		return "", nil
	}

	// Validate file exists
	if _, err := os.Stat(config.ImagePath); os.IsNotExist(err) {
		return "", fmt.Errorf("og-image file not found: %s", config.ImagePath)
	}

	// Get extension and validate format
	ext := strings.ToLower(filepath.Ext(config.ImagePath))
	if !isValidOGImageFormat(ext) {
		return "", fmt.Errorf("unsupported og-image format: %s (use png, jpg, gif, or webp)", ext)
	}

	// Copy to output as og-image.{ext}
	destFilename := "og-image" + ext
	destPath := filepath.Join(outputDir, destFilename)
	if err := copyFile(config.ImagePath, destPath); err != nil {
		return "", fmt.Errorf("failed to copy og-image: %w", err)
	}

	// Build URL
	ogImageURL := "/" + destFilename
	if config.BaseURL != "" {
		ogImageURL = strings.TrimSuffix(config.BaseURL, "/") + "/" + destFilename
	}

	return ogImageURL, nil
}

// isValidOGImageFormat checks if the file extension is a supported OG image format
func isValidOGImageFormat(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return true
	default:
		return false
	}
}

// GetOGImageMimeType returns the MIME type for an OG image file
func GetOGImageMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return ""
	}
}
