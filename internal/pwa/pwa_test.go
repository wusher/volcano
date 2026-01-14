package pwa

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateManifest(t *testing.T) {
	t.Run("basic manifest without icons", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := ManifestConfig{
			SiteTitle: "Test Site",
			BaseURL:   "",
			HasIcons:  false,
		}

		if err := GenerateManifest(tmpDir, config); err != nil {
			t.Fatalf("GenerateManifest() error = %v", err)
		}

		// Read and parse the manifest
		data, err := os.ReadFile(filepath.Join(tmpDir, "manifest.json"))
		if err != nil {
			t.Fatalf("Failed to read manifest.json: %v", err)
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatalf("Failed to parse manifest.json: %v", err)
		}

		if manifest.Name != "Test Site" {
			t.Errorf("Name = %q, want %q", manifest.Name, "Test Site")
		}
		if manifest.ShortName != "Test Site" {
			t.Errorf("ShortName = %q, want %q", manifest.ShortName, "Test Site")
		}
		if manifest.StartURL != "/" {
			t.Errorf("StartURL = %q, want %q", manifest.StartURL, "/")
		}
		if manifest.Scope != "/" {
			t.Errorf("Scope = %q, want %q", manifest.Scope, "/")
		}
		if manifest.Display != "standalone" {
			t.Errorf("Display = %q, want %q", manifest.Display, "standalone")
		}
		if manifest.ThemeColor != DefaultThemeColor {
			t.Errorf("ThemeColor = %q, want %q", manifest.ThemeColor, DefaultThemeColor)
		}
		if manifest.BackgroundColor != "#ffffff" {
			t.Errorf("BackgroundColor = %q, want %q", manifest.BackgroundColor, "#ffffff")
		}
		if len(manifest.Icons) != 0 {
			t.Errorf("Icons should be empty, got %d", len(manifest.Icons))
		}
	})

	t.Run("manifest with custom theme color", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := ManifestConfig{
			SiteTitle:  "Custom Theme",
			ThemeColor: "#ff6600",
			HasIcons:   false,
		}

		if err := GenerateManifest(tmpDir, config); err != nil {
			t.Fatalf("GenerateManifest() error = %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, "manifest.json"))
		if err != nil {
			t.Fatal(err)
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatal(err)
		}

		if manifest.ThemeColor != "#ff6600" {
			t.Errorf("ThemeColor = %q, want %q", manifest.ThemeColor, "#ff6600")
		}
	})

	t.Run("manifest with base URL", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := ManifestConfig{
			SiteTitle: "With Base URL",
			BaseURL:   "/docs",
			HasIcons:  true,
		}

		if err := GenerateManifest(tmpDir, config); err != nil {
			t.Fatalf("GenerateManifest() error = %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, "manifest.json"))
		if err != nil {
			t.Fatal(err)
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatal(err)
		}

		if manifest.StartURL != "/docs/" {
			t.Errorf("StartURL = %q, want %q", manifest.StartURL, "/docs/")
		}
		if manifest.Scope != "/docs/" {
			t.Errorf("Scope = %q, want %q", manifest.Scope, "/docs/")
		}
		if len(manifest.Icons) != 2 {
			t.Errorf("Icons length = %d, want 2", len(manifest.Icons))
		}
		if manifest.Icons[0].Src != "/docs/icon-192.png" {
			t.Errorf("Icon 0 Src = %q, want %q", manifest.Icons[0].Src, "/docs/icon-192.png")
		}
	})

	t.Run("long site title truncated for short name", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := ManifestConfig{
			SiteTitle: "This Is A Very Long Site Title That Should Be Truncated",
			HasIcons:  false,
		}

		if err := GenerateManifest(tmpDir, config); err != nil {
			t.Fatalf("GenerateManifest() error = %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, "manifest.json"))
		if err != nil {
			t.Fatal(err)
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatal(err)
		}

		if len(manifest.ShortName) > 12 {
			t.Errorf("ShortName length = %d, want <= 12", len(manifest.ShortName))
		}
		if manifest.ShortName != "This Is A Ve" {
			t.Errorf("ShortName = %q, want %q", manifest.ShortName, "This Is A Ve")
		}
	})
}

func TestGetManifestLinkTag(t *testing.T) {
	t.Run("without base URL", func(t *testing.T) {
		tag := GetManifestLinkTag("")
		expected := `<link rel="manifest" href="/manifest.json">`
		if tag != expected {
			t.Errorf("GetManifestLinkTag() = %q, want %q", tag, expected)
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		tag := GetManifestLinkTag("/docs")
		expected := `<link rel="manifest" href="/docs/manifest.json">`
		if tag != expected {
			t.Errorf("GetManifestLinkTag() = %q, want %q", tag, expected)
		}
	})
}

func TestGenerateServiceWorker(t *testing.T) {
	t.Run("basic service worker", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := ServiceWorkerConfig{
			BaseURL:   "",
			PageURLs:  []string{"/", "/about/", "/contact/"},
			AssetURLs: []string{"/styles.css", "/app.js"},
		}

		if err := GenerateServiceWorker(tmpDir, config); err != nil {
			t.Fatalf("GenerateServiceWorker() error = %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, "sw.js"))
		if err != nil {
			t.Fatalf("Failed to read sw.js: %v", err)
		}

		content := string(data)

		// Check that it contains expected content
		if !strings.Contains(content, "volcano-cache-") {
			t.Error("Service worker should contain cache name")
		}
		if !strings.Contains(content, `"/"`) {
			t.Error("Service worker should contain root URL")
		}
		if !strings.Contains(content, `"/about/"`) {
			t.Error("Service worker should contain /about/ URL")
		}
		if !strings.Contains(content, `"/styles.css"`) {
			t.Error("Service worker should contain CSS URL")
		}
		if !strings.Contains(content, "self.addEventListener('install'") {
			t.Error("Service worker should contain install event handler")
		}
		if !strings.Contains(content, "self.addEventListener('activate'") {
			t.Error("Service worker should contain activate event handler")
		}
		if !strings.Contains(content, "self.addEventListener('fetch'") {
			t.Error("Service worker should contain fetch event handler")
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		tmpDir := t.TempDir()

		config := ServiceWorkerConfig{
			BaseURL:   "/docs",
			PageURLs:  []string{"/docs/", "/docs/guide/"},
			AssetURLs: []string{"/docs/styles.css"},
		}

		if err := GenerateServiceWorker(tmpDir, config); err != nil {
			t.Fatalf("GenerateServiceWorker() error = %v", err)
		}

		data, err := os.ReadFile(filepath.Join(tmpDir, "sw.js"))
		if err != nil {
			t.Fatal(err)
		}

		content := string(data)
		if !strings.Contains(content, `"/docs/"`) {
			t.Error("Service worker should contain base URL")
		}
	})
}

func TestGetServiceWorkerRegistration(t *testing.T) {
	t.Run("without base URL", func(t *testing.T) {
		js := GetServiceWorkerRegistration("")
		if !strings.Contains(js, "'/sw.js'") {
			t.Error("Registration should use /sw.js")
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		js := GetServiceWorkerRegistration("/docs")
		if !strings.Contains(js, "'/docs/sw.js'") {
			t.Error("Registration should use /docs/sw.js")
		}
	})
}

func TestGenerateIcons(t *testing.T) {
	t.Run("no favicon provided", func(t *testing.T) {
		tmpDir := t.TempDir()

		result, err := GenerateIcons("", tmpDir)
		if err != nil {
			t.Fatalf("GenerateIcons() error = %v", err)
		}

		if result.Generated {
			t.Error("Generated should be false when no favicon provided")
		}
		if len(result.Paths) != 0 {
			t.Error("Paths should be empty")
		}
		if result.Warning != "" {
			t.Errorf("Warning should be empty, got %q", result.Warning)
		}
	})

	t.Run("SVG favicon skipped", func(t *testing.T) {
		tmpDir := t.TempDir()
		faviconPath := filepath.Join(tmpDir, "favicon.svg")
		if err := os.WriteFile(faviconPath, []byte("<svg></svg>"), 0644); err != nil {
			t.Fatal(err)
		}

		result, err := GenerateIcons(faviconPath, tmpDir)
		if err != nil {
			t.Fatalf("GenerateIcons() error = %v", err)
		}

		if result.Generated {
			t.Error("Generated should be false for SVG")
		}
		if result.Warning == "" {
			t.Error("Warning should be set for SVG")
		}
		if !strings.Contains(result.Warning, "SVG") {
			t.Errorf("Warning should mention SVG: %q", result.Warning)
		}
	})

	t.Run("ICO favicon skipped", func(t *testing.T) {
		tmpDir := t.TempDir()
		faviconPath := filepath.Join(tmpDir, "favicon.ico")
		if err := os.WriteFile(faviconPath, []byte{0, 0, 1, 0}, 0644); err != nil {
			t.Fatal(err)
		}

		result, err := GenerateIcons(faviconPath, tmpDir)
		if err != nil {
			t.Fatalf("GenerateIcons() error = %v", err)
		}

		if result.Generated {
			t.Error("Generated should be false for ICO")
		}
		if !strings.Contains(result.Warning, "ICO") {
			t.Errorf("Warning should mention ICO: %q", result.Warning)
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		tmpDir := t.TempDir()
		faviconPath := filepath.Join(tmpDir, "favicon.bmp")
		if err := os.WriteFile(faviconPath, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		result, err := GenerateIcons(faviconPath, tmpDir)
		if err != nil {
			t.Fatalf("GenerateIcons() error = %v", err)
		}

		if result.Generated {
			t.Error("Generated should be false for unsupported format")
		}
		if !strings.Contains(result.Warning, "Unsupported") {
			t.Errorf("Warning should mention unsupported: %q", result.Warning)
		}
	})

	t.Run("PNG favicon generates icons", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a test PNG image (512x512)
		faviconPath := filepath.Join(tmpDir, "favicon.png")
		img := image.NewRGBA(image.Rect(0, 0, 512, 512))
		f, err := os.Create(faviconPath)
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, img); err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		_ = f.Close()

		result, err := GenerateIcons(faviconPath, tmpDir)
		if err != nil {
			t.Fatalf("GenerateIcons() error = %v", err)
		}

		if !result.Generated {
			t.Error("Generated should be true for PNG")
		}
		if len(result.Paths) != 2 {
			t.Errorf("Paths length = %d, want 2", len(result.Paths))
		}

		// Verify icon files exist
		matches, _ := filepath.Glob(filepath.Join(tmpDir, "icon-*.png"))
		if len(matches) != 2 {
			t.Errorf("Expected 2 icon files, got %d", len(matches))
		}

		// No warning for 512x512 image
		if result.Warning != "" {
			t.Errorf("Warning should be empty for 512x512 image, got %q", result.Warning)
		}
	})

	t.Run("small PNG favicon generates warning", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create a small test PNG image (100x100)
		faviconPath := filepath.Join(tmpDir, "small-favicon.png")
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		f, err := os.Create(faviconPath)
		if err != nil {
			t.Fatal(err)
		}
		if err := png.Encode(f, img); err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		_ = f.Close()

		result, err := GenerateIcons(faviconPath, tmpDir)
		if err != nil {
			t.Fatalf("GenerateIcons() error = %v", err)
		}

		if !result.Generated {
			t.Error("Generated should be true for PNG")
		}

		// Check for warning about small image
		if result.Warning == "" {
			t.Error("Should warn about image smaller than 512px")
		}
		if !strings.Contains(result.Warning, "100x100") {
			t.Errorf("Warning should mention image size: %q", result.Warning)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := GenerateIcons(filepath.Join(tmpDir, "nonexistent.png"), tmpDir)
		if err == nil {
			t.Error("GenerateIcons() should return error for non-existent file")
		}
	})
}

func TestGetIconURLs(t *testing.T) {
	t.Run("without base URL", func(t *testing.T) {
		urls := GetIconURLs("")
		if len(urls) != 2 {
			t.Errorf("URLs length = %d, want 2", len(urls))
		}
		if urls[0] != "/icon-192.png" {
			t.Errorf("URL[0] = %q, want %q", urls[0], "/icon-192.png")
		}
		if urls[1] != "/icon-512.png" {
			t.Errorf("URL[1] = %q, want %q", urls[1], "/icon-512.png")
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		urls := GetIconURLs("/docs")
		if urls[0] != "/docs/icon-192.png" {
			t.Errorf("URL[0] = %q, want %q", urls[0], "/docs/icon-192.png")
		}
		if urls[1] != "/docs/icon-512.png" {
			t.Errorf("URL[1] = %q, want %q", urls[1], "/docs/icon-512.png")
		}
	})
}

func TestGenerateIconBytes(t *testing.T) {
	// Create a test PNG image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	pngData := buf.Bytes()

	t.Run("PNG image", func(t *testing.T) {
		result, err := GenerateIconBytes(pngData, ".png", 64)
		if err != nil {
			t.Fatalf("GenerateIconBytes() error = %v", err)
		}
		if len(result) == 0 {
			t.Error("GenerateIconBytes() returned empty result")
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		_, err := GenerateIconBytes([]byte("not an image"), ".bmp", 64)
		if err == nil {
			t.Error("GenerateIconBytes() should return error for unsupported format")
		}
		if !strings.Contains(err.Error(), "unsupported format") {
			t.Errorf("error should mention unsupported format, got: %v", err)
		}
	})

	t.Run("invalid image data", func(t *testing.T) {
		_, err := GenerateIconBytes([]byte("not a valid png"), ".png", 64)
		if err == nil {
			t.Error("GenerateIconBytes() should return error for invalid image data")
		}
	})
}

func TestBuildManifest(t *testing.T) {
	t.Run("basic manifest", func(t *testing.T) {
		config := ManifestConfig{
			SiteTitle:   "Test Site",
			Description: "A test site",
			ThemeColor:  "#ff0000",
			BaseURL:     "",
			HasIcons:    false,
		}

		data := BuildManifest(config)
		if len(data) == 0 {
			t.Error("BuildManifest() returned empty result")
		}

		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatalf("Failed to parse manifest JSON: %v", err)
		}

		if manifest.Name != "Test Site" {
			t.Errorf("Name = %q, want %q", manifest.Name, "Test Site")
		}
		if manifest.ThemeColor != "#ff0000" {
			t.Errorf("ThemeColor = %q, want %q", manifest.ThemeColor, "#ff0000")
		}
	})

	t.Run("manifest with base URL and icons", func(t *testing.T) {
		config := ManifestConfig{
			SiteTitle: "Test Site",
			BaseURL:   "/docs",
			HasIcons:  true,
		}

		data := BuildManifest(config)
		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatalf("Failed to parse manifest JSON: %v", err)
		}

		if manifest.StartURL != "/docs/" {
			t.Errorf("StartURL = %q, want %q", manifest.StartURL, "/docs/")
		}
		if manifest.Scope != "/docs/" {
			t.Errorf("Scope = %q, want %q", manifest.Scope, "/docs/")
		}
		if len(manifest.Icons) != 2 {
			t.Errorf("Icons length = %d, want 2", len(manifest.Icons))
		}
		if manifest.Icons[0].Src != "/docs/icon-192.png" {
			t.Errorf("Icon[0].Src = %q, want %q", manifest.Icons[0].Src, "/docs/icon-192.png")
		}
	})

	t.Run("default theme color", func(t *testing.T) {
		config := ManifestConfig{
			SiteTitle: "Test Site",
		}

		data := BuildManifest(config)
		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatalf("Failed to parse manifest JSON: %v", err)
		}

		if manifest.ThemeColor != DefaultThemeColor {
			t.Errorf("ThemeColor = %q, want default %q", manifest.ThemeColor, DefaultThemeColor)
		}
	})

	t.Run("long title truncated", func(t *testing.T) {
		config := ManifestConfig{
			SiteTitle: "This Is A Very Long Site Title",
		}

		data := BuildManifest(config)
		var manifest Manifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			t.Fatalf("Failed to parse manifest JSON: %v", err)
		}

		if len(manifest.ShortName) > 12 {
			t.Errorf("ShortName length = %d, want <= 12", len(manifest.ShortName))
		}
	})
}

func TestBuildServiceWorker(t *testing.T) {
	t.Run("basic service worker", func(t *testing.T) {
		config := ServiceWorkerConfig{
			BaseURL:   "",
			PageURLs:  []string{"/", "/about/"},
			AssetURLs: []string{"/styles.css"},
		}

		result := BuildServiceWorker(config)
		if result == "" {
			t.Error("BuildServiceWorker() returned empty result")
		}

		if !strings.Contains(result, "volcano-cache-") {
			t.Error("service worker should contain cache name")
		}
		if !strings.Contains(result, `"/"`) {
			t.Error("service worker should contain root URL")
		}
		if !strings.Contains(result, `"/about/"`) {
			t.Error("service worker should contain about URL")
		}
		if !strings.Contains(result, `"/styles.css"`) {
			t.Error("service worker should contain CSS URL")
		}
	})

	t.Run("deterministic cache version", func(t *testing.T) {
		config := ServiceWorkerConfig{
			PageURLs:  []string{"/a/", "/b/"},
			AssetURLs: []string{"/c.css"},
		}

		result1 := BuildServiceWorker(config)
		result2 := BuildServiceWorker(config)

		if result1 != result2 {
			t.Error("BuildServiceWorker() should return identical results for same input")
		}
	})
}
