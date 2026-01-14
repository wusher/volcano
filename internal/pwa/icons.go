package pwa

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

// IconSizes defines the sizes for PWA icons.
var IconSizes = []int{192, 512}

// IconResult holds the result of icon generation.
type IconResult struct {
	Generated bool     // Whether icons were successfully generated
	Paths     []string // Paths to generated icon files
	Warning   string   // Warning if source too small or unsupported format
}

// GenerateIcons creates PWA icons from a favicon source.
func GenerateIcons(faviconPath, outputDir string) (*IconResult, error) {
	result := &IconResult{}

	// If no favicon provided, skip icon generation
	if faviconPath == "" {
		return result, nil
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(faviconPath))

	// Skip unsupported formats
	switch ext {
	case ".svg":
		result.Warning = "SVG favicon cannot be resized for PWA icons (requires rasterization). Skipping icon generation."
		return result, nil
	case ".ico":
		result.Warning = "ICO favicon cannot be resized for PWA icons (complex container format). Skipping icon generation."
		return result, nil
	case ".png", ".jpg", ".jpeg", ".gif":
		// Supported formats
	default:
		result.Warning = fmt.Sprintf("Unsupported favicon format %q for PWA icons. Use PNG, JPG, or GIF.", ext)
		return result, nil
	}

	// Open and decode the source image
	f, err := os.Open(faviconPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open favicon: %w", err)
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
		return nil, fmt.Errorf("failed to decode favicon: %w", err)
	}

	// Check source size
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	if srcWidth < 512 || srcHeight < 512 {
		result.Warning = fmt.Sprintf("Favicon is %dx%d, smaller than recommended 512x512. PWA icons may appear blurry.", srcWidth, srcHeight)
	}

	// Generate icons for each size
	for _, size := range IconSizes {
		iconPath := filepath.Join(outputDir, fmt.Sprintf("icon-%d.png", size))
		if err := resizeAndSave(src, size, iconPath); err != nil {
			return nil, fmt.Errorf("failed to generate %dx%d icon: %w", size, size, err)
		}
		result.Paths = append(result.Paths, iconPath)
	}

	result.Generated = true
	return result, nil
}

// resizeAndSave resizes an image to the given size and saves it as PNG.
func resizeAndSave(src image.Image, size int, outputPath string) error {
	// Create destination image
	dst := image.NewRGBA(image.Rect(0, 0, size, size))

	// Use high-quality CatmullRom resampling
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	// Encode as PNG
	return png.Encode(f, dst)
}

// GetIconURLs returns the URLs for PWA icons.
func GetIconURLs(baseURL string) []string {
	prefix := "/"
	if baseURL != "" {
		prefix = baseURL + "/"
	}

	urls := make([]string, len(IconSizes))
	for i, size := range IconSizes {
		urls[i] = fmt.Sprintf("%sicon-%d.png", prefix, size)
	}
	return urls
}
