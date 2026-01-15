package assets

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"favicon.ico", "image/x-icon"},
		{"favicon.png", "image/png"},
		{"favicon.svg", "image/svg+xml"},
		{"favicon.gif", "image/gif"},
		{"favicon.ICO", "image/x-icon"},
		{"favicon.PNG", "image/png"},
		{"favicon.txt", ""},
		{"favicon.jpg", ""},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			result := GetFaviconMimeType(tc.filename)
			if result != tc.expected {
				t.Errorf("GetFaviconMimeType(%q) = %q, want %q", tc.filename, result, tc.expected)
			}
		})
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := "test content"
	if err := os.WriteFile(srcPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Copy file
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile error: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("copied content = %q, want %q", string(data), content)
	}
}

func TestCopyFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	err := copyFile(filepath.Join(tmpDir, "nonexistent"), filepath.Join(tmpDir, "dest"))
	if err == nil {
		t.Error("expected error for nonexistent source")
	}
}

func TestProcessFavicon(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test favicon
	faviconPath := filepath.Join(tmpDir, "favicon.ico")
	if err := os.WriteFile(faviconPath, []byte("fake ico"), 0644); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{IconPath: faviconPath}
	links, err := ProcessFavicon(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}

	if links[0].Rel != "icon" {
		t.Errorf("Rel = %q, want %q", links[0].Rel, "icon")
	}
	if links[0].Type != "image/x-icon" {
		t.Errorf("Type = %q, want %q", links[0].Type, "image/x-icon")
	}
	if links[0].Href != "/favicon.ico" {
		t.Errorf("Href = %q, want %q", links[0].Href, "/favicon.ico")
	}

	// Check file was copied
	if _, err := os.Stat(filepath.Join(outputDir, "favicon.ico")); os.IsNotExist(err) {
		t.Error("favicon.ico not copied to output")
	}
}

func TestProcessFaviconNotFound(t *testing.T) {
	config := FaviconConfig{IconPath: "/nonexistent/favicon.ico"}
	_, err := ProcessFavicon(config, "/tmp")
	if err == nil {
		t.Error("expected error for nonexistent favicon")
	}
}

func TestProcessFaviconUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with unsupported extension
	badPath := filepath.Join(tmpDir, "favicon.txt")
	if err := os.WriteFile(badPath, []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{IconPath: badPath}
	_, err := ProcessFavicon(config, tmpDir)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestProcessFaviconEmpty(t *testing.T) {
	config := FaviconConfig{}
	links, err := ProcessFavicon(config, "/tmp")
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 links for empty config, got %d", len(links))
	}
}

func TestProcessFaviconWithAppleTouchIcon(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	iconPath := filepath.Join(tmpDir, "favicon.png")
	if err := os.WriteFile(iconPath, []byte("png"), 0644); err != nil {
		t.Fatal(err)
	}

	applePath := filepath.Join(tmpDir, "apple-touch-icon.png")
	if err := os.WriteFile(applePath, []byte("png"), 0644); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{
		IconPath:       iconPath,
		AppleTouchIcon: applePath,
	}
	links, err := ProcessFavicon(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}

	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}

	if links[1].Rel != "apple-touch-icon" {
		t.Errorf("Rel = %q, want %q", links[1].Rel, "apple-touch-icon")
	}
	if links[1].Sizes != "180x180" {
		t.Errorf("Sizes = %q, want %q", links[1].Sizes, "180x180")
	}
}

func TestRenderFaviconLinks(t *testing.T) {
	links := []FaviconLink{
		{Rel: "icon", Type: "image/x-icon", Href: "/favicon.ico"},
		{Rel: "apple-touch-icon", Type: "image/png", Sizes: "180x180", Href: "/apple-touch-icon.png"},
	}

	html := string(RenderFaviconLinks(links))

	if !strings.Contains(html, `rel="icon"`) {
		t.Error("missing icon rel")
	}
	if !strings.Contains(html, `type="image/x-icon"`) {
		t.Error("missing icon type")
	}
	if !strings.Contains(html, `href="/favicon.ico"`) {
		t.Error("missing icon href")
	}
	if !strings.Contains(html, `rel="apple-touch-icon"`) {
		t.Error("missing apple-touch-icon rel")
	}
	if !strings.Contains(html, `sizes="180x180"`) {
		t.Error("missing sizes")
	}
}

func TestRenderFaviconLinksEmpty(t *testing.T) {
	html := RenderFaviconLinks(nil)
	if html != "" {
		t.Errorf("expected empty string, got %q", html)
	}
}

func TestGenerateAppleIcons_SVGNotSupported(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test SVG file
	svgPath := filepath.Join(tmpDir, "icon.svg")
	if err := os.WriteFile(svgPath, []byte("<svg></svg>"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := generateAppleIcons(svgPath, tmpDir)
	if err == nil {
		t.Error("generateAppleIcons should return error for SVG format")
	}
}

func TestGenerateAppleIcons_ICONotSupported(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test ICO file
	icoPath := filepath.Join(tmpDir, "icon.ico")
	if err := os.WriteFile(icoPath, []byte("fake ico data"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := generateAppleIcons(icoPath, tmpDir)
	if err == nil {
		t.Error("generateAppleIcons should return error for ICO format")
	}
}

func TestGenerateAppleIcons_UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with unknown extension
	unknownPath := filepath.Join(tmpDir, "icon.xyz")
	if err := os.WriteFile(unknownPath, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := generateAppleIcons(unknownPath, tmpDir)
	if err == nil {
		t.Error("generateAppleIcons should return error for unsupported format")
	}
}

func TestGenerateAppleIcons_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := generateAppleIcons(filepath.Join(tmpDir, "nonexistent.png"), tmpDir)
	if err == nil {
		t.Error("generateAppleIcons should return error for nonexistent file")
	}
}

func TestProcessFaviconWithPNG_GeneratesAppleIcons(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a valid PNG file
	pngPath := filepath.Join(tmpDir, "logo.png")
	if err := createTestPNG(pngPath, 100, 100); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{IconPath: pngPath}
	links, err := ProcessFavicon(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}

	// Should have original favicon + apple-touch-icon
	if len(links) < 2 {
		t.Errorf("expected at least 2 links (favicon + apple icon), got %d", len(links))
	}

	// Check apple-touch-icon.png was created
	applePath := filepath.Join(outputDir, "apple-touch-icon.png")
	if _, err := os.Stat(applePath); os.IsNotExist(err) {
		t.Error("apple-touch-icon.png should be generated")
	}

	// Check apple-touch-icon-precomposed.png was created
	precomposedPath := filepath.Join(outputDir, "apple-touch-icon-precomposed.png")
	if _, err := os.Stat(precomposedPath); os.IsNotExist(err) {
		t.Error("apple-touch-icon-precomposed.png should be generated")
	}

	// Check favicon.ico was created
	faviconIcoPath := filepath.Join(outputDir, "favicon.ico")
	if _, err := os.Stat(faviconIcoPath); os.IsNotExist(err) {
		t.Error("favicon.ico should be generated from PNG source")
	}
}

func TestGenerateAppleIcons_WithJPEG(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid JPEG file
	jpegPath := filepath.Join(tmpDir, "logo.jpg")
	if err := createTestJPEG(jpegPath, 100, 100); err != nil {
		t.Fatal(err)
	}

	// generateAppleIcons supports JPEG
	links, err := generateAppleIcons(jpegPath, tmpDir)
	if err != nil {
		t.Fatalf("generateAppleIcons error: %v", err)
	}

	// Should have apple-touch-icon generated
	hasAppleIcon := false
	for _, link := range links {
		if link.Rel == "apple-touch-icon" {
			hasAppleIcon = true
			break
		}
	}
	if !hasAppleIcon {
		t.Error("Should generate apple-touch-icon from JPEG")
	}
}

func TestGenerateAppleIcons_WithGIF(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid GIF file
	gifPath := filepath.Join(tmpDir, "logo.gif")
	if err := createTestGIF(gifPath, 100, 100); err != nil {
		t.Fatal(err)
	}

	// generateAppleIcons supports GIF
	links, err := generateAppleIcons(gifPath, tmpDir)
	if err != nil {
		t.Fatalf("generateAppleIcons error: %v", err)
	}

	// Should have apple-touch-icon generated
	hasAppleIcon := false
	for _, link := range links {
		if link.Rel == "apple-touch-icon" {
			hasAppleIcon = true
			break
		}
	}
	if !hasAppleIcon {
		t.Error("Should generate apple-touch-icon from GIF")
	}
}

func TestResizeAndSavePNG(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test image
	srcPath := filepath.Join(tmpDir, "source.png")
	if err := createTestPNG(srcPath, 200, 200); err != nil {
		t.Fatal(err)
	}

	// Load the source image
	f, err := os.Open(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	src, err := pngDecode(f)
	_ = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Resize to 64x64
	dstPath := filepath.Join(tmpDir, "resized.png")
	if err := resizeAndSavePNG(src, 64, dstPath); err != nil {
		t.Fatalf("resizeAndSavePNG error: %v", err)
	}

	// Verify output exists
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		t.Error("resized PNG should be created")
	}

	// Verify dimensions
	resized, err := os.Open(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resized.Close() }()

	img, err := pngDecode(resized)
	if err != nil {
		t.Fatal(err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		t.Errorf("resized image should be 64x64, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestProcessFaviconWithAppleTouchIconOverride(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a valid PNG favicon
	pngPath := filepath.Join(tmpDir, "logo.png")
	if err := createTestPNG(pngPath, 100, 100); err != nil {
		t.Fatal(err)
	}

	// Create a separate apple touch icon
	applePath := filepath.Join(tmpDir, "apple-custom.png")
	if err := createTestPNG(applePath, 180, 180); err != nil {
		t.Fatal(err)
	}

	config := FaviconConfig{
		IconPath:       pngPath,
		AppleTouchIcon: applePath,
	}
	links, err := ProcessFavicon(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessFavicon error: %v", err)
	}

	// Should have apple-touch-icon with custom href
	hasAppleIcon := false
	for _, link := range links {
		if link.Rel == "apple-touch-icon" {
			hasAppleIcon = true
			if link.Href != "/apple-custom.png" {
				t.Errorf("Expected custom apple icon href, got %s", link.Href)
			}
			break
		}
	}
	if !hasAppleIcon {
		t.Error("Should have apple-touch-icon link")
	}
}

// Helper functions to create test images

func createTestPNG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with a color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, image.Black)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return png.Encode(f, img)
}

func createTestJPEG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, image.Black)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return jpeg.Encode(f, img, nil)
}

func createTestGIF(path string, width, height int) error {
	// Use a simple palette with black and white
	palette := color.Palette{color.Black, color.White}
	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return gif.Encode(f, img, nil)
}

// Wrapper functions to use real image package
func pngDecode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}
