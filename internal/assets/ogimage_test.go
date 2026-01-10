package assets

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidOGImageFormat(t *testing.T) {
	// Note: isValidOGImageFormat expects lowercase extensions
	// ProcessOGImage calls strings.ToLower before calling this function
	tests := []struct {
		ext      string
		expected bool
	}{
		{".png", true},
		{".jpg", true},
		{".jpeg", true},
		{".gif", true},
		{".webp", true},
		{".ico", false},
		{".svg", false},
		{".txt", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.ext, func(t *testing.T) {
			result := isValidOGImageFormat(tc.ext)
			if result != tc.expected {
				t.Errorf("isValidOGImageFormat(%q) = %v, want %v", tc.ext, result, tc.expected)
			}
		})
	}
}

func TestGetOGImageMimeType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"image.png", "image/png"},
		{"image.jpg", "image/jpeg"},
		{"image.jpeg", "image/jpeg"},
		{"image.gif", "image/gif"},
		{"image.webp", "image/webp"},
		{"image.PNG", "image/png"},
		{"image.JPG", "image/jpeg"},
		{"image.ico", ""},
		{"image.txt", ""},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			result := GetOGImageMimeType(tc.filename)
			if result != tc.expected {
				t.Errorf("GetOGImageMimeType(%q) = %q, want %q", tc.filename, result, tc.expected)
			}
		})
	}
}

func TestProcessOGImage(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test OG image
	ogImagePath := filepath.Join(tmpDir, "social.png")
	if err := os.WriteFile(ogImagePath, []byte("fake png"), 0644); err != nil {
		t.Fatal(err)
	}

	config := OGImageConfig{
		ImagePath: ogImagePath,
		BaseURL:   "https://example.com",
	}
	url, err := ProcessOGImage(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessOGImage error: %v", err)
	}

	expectedURL := "https://example.com/og-image.png"
	if url != expectedURL {
		t.Errorf("URL = %q, want %q", url, expectedURL)
	}

	// Check file was copied with standardized name
	if _, err := os.Stat(filepath.Join(outputDir, "og-image.png")); os.IsNotExist(err) {
		t.Error("og-image.png not copied to output")
	}
}

func TestProcessOGImageWithoutBaseURL(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a test OG image
	ogImagePath := filepath.Join(tmpDir, "social.jpg")
	if err := os.WriteFile(ogImagePath, []byte("fake jpg"), 0644); err != nil {
		t.Fatal(err)
	}

	config := OGImageConfig{
		ImagePath: ogImagePath,
		BaseURL:   "", // No base URL
	}
	url, err := ProcessOGImage(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessOGImage error: %v", err)
	}

	expectedURL := "/og-image.jpg"
	if url != expectedURL {
		t.Errorf("URL = %q, want %q", url, expectedURL)
	}
}

func TestProcessOGImageEmpty(t *testing.T) {
	config := OGImageConfig{}
	url, err := ProcessOGImage(config, "/tmp")
	if err != nil {
		t.Fatalf("ProcessOGImage error: %v", err)
	}
	if url != "" {
		t.Errorf("expected empty URL for empty config, got %q", url)
	}
}

func TestProcessOGImageNotFound(t *testing.T) {
	config := OGImageConfig{ImagePath: "/nonexistent/image.png"}
	_, err := ProcessOGImage(config, "/tmp")
	if err == nil {
		t.Error("expected error for nonexistent image")
	}
}

func TestProcessOGImageUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file with unsupported extension
	badPath := filepath.Join(tmpDir, "image.ico")
	if err := os.WriteFile(badPath, []byte("fake ico"), 0644); err != nil {
		t.Fatal(err)
	}

	config := OGImageConfig{ImagePath: badPath}
	_, err := ProcessOGImage(config, tmpDir)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestProcessOGImageWithTrailingSlash(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatal(err)
	}

	ogImagePath := filepath.Join(tmpDir, "social.webp")
	if err := os.WriteFile(ogImagePath, []byte("fake webp"), 0644); err != nil {
		t.Fatal(err)
	}

	config := OGImageConfig{
		ImagePath: ogImagePath,
		BaseURL:   "https://example.com/", // With trailing slash
	}
	url, err := ProcessOGImage(config, outputDir)
	if err != nil {
		t.Fatalf("ProcessOGImage error: %v", err)
	}

	expectedURL := "https://example.com/og-image.webp"
	if url != expectedURL {
		t.Errorf("URL = %q, want %q", url, expectedURL)
	}
}
