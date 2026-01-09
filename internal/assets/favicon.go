// Package assets handles static asset processing like favicons.
package assets

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FaviconConfig holds configuration for favicon handling
type FaviconConfig struct {
	IconPath       string // Path to favicon file (ico, png, or svg)
	AppleTouchIcon string // Path to Apple touch icon (optional)
}

// FaviconLink represents a favicon link tag
type FaviconLink struct {
	Rel   string // "icon" or "apple-touch-icon"
	Type  string // MIME type
	Sizes string // Size attribute (optional)
	Href  string // URL path
}

// ProcessFavicon copies the favicon to the output directory and returns link tags
func ProcessFavicon(config FaviconConfig, outputDir string) ([]FaviconLink, error) {
	var links []FaviconLink

	// Process main favicon
	if config.IconPath != "" {
		link, err := processSingleFavicon(config.IconPath, outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to process favicon: %w", err)
		}
		links = append(links, link)
	}

	// Process Apple touch icon
	if config.AppleTouchIcon != "" {
		link, err := processSingleFavicon(config.AppleTouchIcon, outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to process apple-touch-icon: %w", err)
		}
		link.Rel = "apple-touch-icon"
		link.Sizes = "180x180"
		links = append(links, link)
	}

	return links, nil
}

// processSingleFavicon copies a favicon file and returns its link
func processSingleFavicon(sourcePath string, outputDir string) (FaviconLink, error) {
	// Validate file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return FaviconLink{}, fmt.Errorf("favicon file not found: %s", sourcePath)
	}

	// Get filename and MIME type
	filename := filepath.Base(sourcePath)
	mimeType := getMimeType(filename)
	if mimeType == "" {
		return FaviconLink{}, fmt.Errorf("unsupported favicon format: %s", filename)
	}

	// Copy file to output directory
	destPath := filepath.Join(outputDir, filename)
	if err := copyFile(sourcePath, destPath); err != nil {
		return FaviconLink{}, err
	}

	return FaviconLink{
		Rel:  "icon",
		Type: mimeType,
		Href: "/" + filename,
	}, nil
}

// getMimeType returns the MIME type for a favicon file
func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".ico":
		return "image/x-icon"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	case ".gif":
		return "image/gif"
	default:
		return ""
	}
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// RenderFaviconLinks generates HTML link tags for favicons
func RenderFaviconLinks(links []FaviconLink) template.HTML {
	if len(links) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, link := range links {
		sb.WriteString(`<link rel="`)
		sb.WriteString(link.Rel)
		sb.WriteString(`"`)
		if link.Type != "" {
			sb.WriteString(` type="`)
			sb.WriteString(link.Type)
			sb.WriteString(`"`)
		}
		if link.Sizes != "" {
			sb.WriteString(` sizes="`)
			sb.WriteString(link.Sizes)
			sb.WriteString(`"`)
		}
		sb.WriteString(` href="`)
		sb.WriteString(link.Href)
		sb.WriteString(`">`)
		sb.WriteString("\n")
	}
	return template.HTML(sb.String())
}
