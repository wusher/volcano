// Package assets handles static asset processing like favicons.
package assets

import (
	"fmt"
	"html/template"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
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

// ProcessFavicon copies the favicon to the output directory and returns link tags.
// It also generates Apple touch icons and favicon.ico for broad browser compatibility.
func ProcessFavicon(config FaviconConfig, outputDir string) ([]FaviconLink, error) {
	var links []FaviconLink

	if config.IconPath == "" {
		return links, nil
	}

	// Process main favicon
	link, err := processSingleFavicon(config.IconPath, outputDir)
	if err != nil {
		return nil, fmt.Errorf("failed to process favicon: %w", err)
	}
	links = append(links, link)

	// Generate Apple touch icons and favicon.ico from the source
	// Non-fatal - just skip Apple icons if we can't generate them (e.g., for SVG or ICO sources)
	if appleLinks, err := generateAppleIcons(config.IconPath, outputDir); err == nil {
		links = append(links, appleLinks...)
	}

	// Process explicit Apple touch icon if provided (overrides generated one)
	if config.AppleTouchIcon != "" {
		applLink, err := processSingleFavicon(config.AppleTouchIcon, outputDir)
		if err != nil {
			return nil, fmt.Errorf("failed to process apple-touch-icon: %w", err)
		}
		applLink.Rel = "apple-touch-icon"
		applLink.Sizes = "180x180"
		// Replace any existing apple-touch-icon link
		for i, l := range links {
			if l.Rel == "apple-touch-icon" {
				links[i] = applLink
				applLink.Href = "" // Mark as handled
				break
			}
		}
		if applLink.Href != "" {
			links = append(links, applLink)
		}
	}

	return links, nil
}

// generateAppleIcons creates apple-touch-icon.png, apple-touch-icon-precomposed.png,
// and favicon.ico from a source image for maximum browser compatibility.
func generateAppleIcons(sourcePath, outputDir string) ([]FaviconLink, error) {
	var links []FaviconLink

	ext := strings.ToLower(filepath.Ext(sourcePath))

	// Skip formats we can't resize
	switch ext {
	case ".svg", ".ico":
		return nil, fmt.Errorf("cannot resize %s format", ext)
	case ".png", ".jpg", ".jpeg", ".gif":
		// Supported
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}

	// Open and decode the source image
	f, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var src image.Image
	switch ext {
	case ".png":
		src, err = png.Decode(f)
	case ".jpg", ".jpeg":
		src, err = jpeg.Decode(f)
	case ".gif":
		src, err = gif.Decode(f)
	}
	if err != nil {
		return nil, err
	}

	// Generate apple-touch-icon.png (180x180)
	appleTouchPath := filepath.Join(outputDir, "apple-touch-icon.png")
	if err := resizeAndSavePNG(src, 180, appleTouchPath); err != nil {
		return nil, err
	}
	links = append(links, FaviconLink{
		Rel:   "apple-touch-icon",
		Type:  "image/png",
		Sizes: "180x180",
		Href:  "/apple-touch-icon.png",
	})

	// Generate apple-touch-icon-precomposed.png (same as above)
	precomposedPath := filepath.Join(outputDir, "apple-touch-icon-precomposed.png")
	if err := copyFile(appleTouchPath, precomposedPath); err != nil {
		return nil, err
	}

	// Generate favicon.ico (32x32 PNG saved as .ico - works in modern browsers)
	faviconPath := filepath.Join(outputDir, "favicon.ico")
	// Check if source is already favicon.ico
	if filepath.Base(sourcePath) != "favicon.ico" {
		if err := resizeAndSavePNG(src, 32, faviconPath); err != nil {
			return nil, err
		}
	}

	return links, nil
}

// resizeAndSavePNG resizes an image to the given square size and saves as PNG.
func resizeAndSavePNG(src image.Image, size int, outputPath string) error {
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return png.Encode(f, dst)
}

// processSingleFavicon copies a favicon file and returns its link
func processSingleFavicon(sourcePath string, outputDir string) (FaviconLink, error) {
	// Validate file exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return FaviconLink{}, fmt.Errorf("favicon file not found: %s", sourcePath)
	}

	// Get filename and MIME type
	filename := filepath.Base(sourcePath)
	mimeType := GetFaviconMimeType(filename)
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

// GetFaviconMimeType returns the MIME type for a favicon file.
// Exported for reuse by dynamic server.
func GetFaviconMimeType(filename string) string {
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
