// Package assets provides asset handling utilities.
package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// HashedAsset represents a static asset with a content-hashed filename.
type HashedAsset struct {
	FileName string // The hashed filename (e.g., "styles.a1b2c3d4.css")
	URLPath  string // The URL path (e.g., "/assets/styles.a1b2c3d4.css")
	Content  string // The raw content
}

// WriteHashedAsset writes content to a file with a content-based hash in the filename.
// Returns the hashed filename and URL path.
// Pattern: {name}.{hash}.{ext} (e.g., "styles.a1b2c3d4.css")
func WriteHashedAsset(outputDir, name, ext, content, baseURL string) (*HashedAsset, error) {
	// Compute SHA256 hash and take first 8 characters
	hash := sha256.Sum256([]byte(content))
	hashStr := hex.EncodeToString(hash[:])[:8]

	// Build filename and paths
	fileName := fmt.Sprintf("%s.%s.%s", name, hashStr, ext)
	assetsDir := filepath.Join(outputDir, "assets")
	filePath := filepath.Join(assetsDir, fileName)

	// Create assets directory
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create assets directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write asset %s: %w", fileName, err)
	}

	// Build URL path with base URL prefix
	urlPath := "/assets/" + fileName
	if baseURL != "" {
		urlPath = baseURL + urlPath
	}

	return &HashedAsset{
		FileName: fileName,
		URLPath:  urlPath,
		Content:  content,
	}, nil
}
