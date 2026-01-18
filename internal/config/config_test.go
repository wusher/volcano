package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	t.Run("valid config file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "valid.json")
		content := `{
			"title": "Test Site",
			"output": "./public",
			"url": "https://example.com",
			"author": "Test Author",
			"theme": "blog",
			"css": "./custom.css",
			"accentColor": "#ff6600",
			"favicon": "./favicon.png",
			"ogImage": "./og.png",
			"topNav": true,
			"breadcrumbs": false,
			"pageNav": true,
			"instantNav": true,
			"inlineAssets": true,
			"pwa": true
		}`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := Load(configPath)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		// Check string values
		if cfg.Title != "Test Site" {
			t.Errorf("Title = %q, want %q", cfg.Title, "Test Site")
		}
		if cfg.Output != "./public" {
			t.Errorf("Output = %q, want %q", cfg.Output, "./public")
		}
		if cfg.URL != "https://example.com" {
			t.Errorf("URL = %q, want %q", cfg.URL, "https://example.com")
		}
		if cfg.Author != "Test Author" {
			t.Errorf("Author = %q, want %q", cfg.Author, "Test Author")
		}
		if cfg.Theme != "blog" {
			t.Errorf("Theme = %q, want %q", cfg.Theme, "blog")
		}
		if cfg.CSS != "./custom.css" {
			t.Errorf("CSS = %q, want %q", cfg.CSS, "./custom.css")
		}
		if cfg.AccentColor != "#ff6600" {
			t.Errorf("AccentColor = %q, want %q", cfg.AccentColor, "#ff6600")
		}
		if cfg.Favicon != "./favicon.png" {
			t.Errorf("Favicon = %q, want %q", cfg.Favicon, "./favicon.png")
		}
		if cfg.OGImage != "./og.png" {
			t.Errorf("OGImage = %q, want %q", cfg.OGImage, "./og.png")
		}

		// Check bool values
		if cfg.TopNav == nil || !*cfg.TopNav {
			t.Error("TopNav should be true")
		}
		if cfg.Breadcrumbs == nil || *cfg.Breadcrumbs {
			t.Error("Breadcrumbs should be false")
		}
		if cfg.PageNav == nil || !*cfg.PageNav {
			t.Error("PageNav should be true")
		}
		if cfg.InstantNav == nil || !*cfg.InstantNav {
			t.Error("InstantNav should be true")
		}
		if cfg.InlineAssets == nil || !*cfg.InlineAssets {
			t.Error("InlineAssets should be true")
		}
		if cfg.PWA == nil || !*cfg.PWA {
			t.Error("PWA should be true")
		}
	})

	t.Run("minimal config file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "minimal.json")
		content := `{"title": "Minimal"}`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := Load(configPath)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Title != "Minimal" {
			t.Errorf("Title = %q, want %q", cfg.Title, "Minimal")
		}
		// All other fields should be empty/nil
		if cfg.Output != "" {
			t.Errorf("Output should be empty, got %q", cfg.Output)
		}
		if cfg.TopNav != nil {
			t.Error("TopNav should be nil")
		}
	})

	t.Run("empty config file", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "empty.json")
		content := `{}`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, err := Load(configPath)
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Title != "" {
			t.Errorf("Title should be empty, got %q", cfg.Title)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "invalid.json")
		content := `{invalid json}`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		_, err := Load(configPath)
		if err == nil {
			t.Error("Load() should return error for invalid JSON")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := Load(filepath.Join(tmpDir, "nonexistent.json"))
		if err == nil {
			t.Error("Load() should return error for non-existent file")
		}
	})
}

func TestDiscover(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("config file exists", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, ConfigFileName)
		if err := os.WriteFile(configPath, []byte(`{}`), 0644); err != nil {
			t.Fatal(err)
		}

		result := Discover(tmpDir)
		if result != configPath {
			t.Errorf("Discover() = %q, want %q", result, configPath)
		}
	})

	t.Run("config file does not exist", func(t *testing.T) {
		emptyDir := t.TempDir()
		result := Discover(emptyDir)
		if result != "" {
			t.Errorf("Discover() = %q, want empty string", result)
		}
	})
}

func TestLoadOrDiscover(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("explicit path exists", func(t *testing.T) {
		configPath := filepath.Join(tmpDir, "explicit.json")
		content := `{"title": "Explicit"}`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, path, err := LoadOrDiscover(configPath, tmpDir)
		if err != nil {
			t.Fatalf("LoadOrDiscover() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadOrDiscover() returned nil config")
		}
		if cfg.Title != "Explicit" {
			t.Errorf("Title = %q, want %q", cfg.Title, "Explicit")
		}
		if path != configPath {
			t.Errorf("path = %q, want %q", path, configPath)
		}
	})

	t.Run("explicit path does not exist", func(t *testing.T) {
		_, _, err := LoadOrDiscover(filepath.Join(tmpDir, "nonexistent.json"), tmpDir)
		if err == nil {
			t.Error("LoadOrDiscover() should return error for non-existent explicit path")
		}
	})

	t.Run("discover config in input dir", func(t *testing.T) {
		inputDir := t.TempDir()
		configPath := filepath.Join(inputDir, ConfigFileName)
		content := `{"title": "Discovered"}`
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		cfg, path, err := LoadOrDiscover("", inputDir)
		if err != nil {
			t.Fatalf("LoadOrDiscover() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadOrDiscover() returned nil config")
		}
		if cfg.Title != "Discovered" {
			t.Errorf("Title = %q, want %q", cfg.Title, "Discovered")
		}
		if path != configPath {
			t.Errorf("path = %q, want %q", path, configPath)
		}
	})

	t.Run("no config found", func(t *testing.T) {
		emptyDir := t.TempDir()
		cfg, path, err := LoadOrDiscover("", emptyDir)
		if err != nil {
			t.Fatalf("LoadOrDiscover() error = %v", err)
		}
		if cfg != nil {
			t.Error("LoadOrDiscover() should return nil config when no file found")
		}
		if path != "" {
			t.Errorf("path should be empty, got %q", path)
		}
	})

	t.Run("discovered config with invalid JSON", func(t *testing.T) {
		inputDir := t.TempDir()
		configPath := filepath.Join(inputDir, ConfigFileName)
		invalidContent := `{invalid json`
		if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
			t.Fatal(err)
		}

		_, path, err := LoadOrDiscover("", inputDir)
		if err == nil {
			t.Error("LoadOrDiscover() should return error for invalid JSON")
		}
		if path != configPath {
			t.Errorf("path = %q, want %q", path, configPath)
		}
	})
}

func TestBoolPtr(t *testing.T) {
	truePtr := BoolPtr(true)
	if truePtr == nil || !*truePtr {
		t.Error("BoolPtr(true) should return pointer to true")
	}

	falsePtr := BoolPtr(false)
	if falsePtr == nil || *falsePtr {
		t.Error("BoolPtr(false) should return pointer to false")
	}
}

func TestGetBool(t *testing.T) {
	trueVal := true
	falseVal := false

	if !GetBool(&trueVal, false) {
		t.Error("GetBool with true pointer should return true")
	}
	if GetBool(&falseVal, true) {
		t.Error("GetBool with false pointer should return false")
	}
	if !GetBool(nil, true) {
		t.Error("GetBool with nil should return default (true)")
	}
	if GetBool(nil, false) {
		t.Error("GetBool with nil should return default (false)")
	}
}

func TestIntPtrAndGetInt(t *testing.T) {
	intPtr := IntPtr(42)
	if intPtr == nil || *intPtr != 42 {
		t.Error("IntPtr(42) should return pointer to 42")
	}

	value := 5
	if GetInt(&value, 10) != 5 {
		t.Error("GetInt with pointer should return value")
	}
	if GetInt(nil, 10) != 10 {
		t.Error("GetInt with nil should return default")
	}
}

func TestDefaultFileConfig(t *testing.T) {
	cfg := DefaultFileConfig()
	if cfg.Output != "./output" {
		t.Errorf("Output = %q, want %q", cfg.Output, "./output")
	}
	if cfg.Port == nil || *cfg.Port != 1776 {
		t.Errorf("Port = %v, want 1776", cfg.Port)
	}
	if cfg.Title != "My Site" {
		t.Errorf("Title = %q, want %q", cfg.Title, "My Site")
	}
	if cfg.Theme != "docs" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "docs")
	}
	if cfg.TopNav == nil || *cfg.TopNav {
		t.Error("TopNav should default to false")
	}
	if cfg.Breadcrumbs == nil || !*cfg.Breadcrumbs {
		t.Error("Breadcrumbs should default to true")
	}
	if cfg.PageNav == nil || *cfg.PageNav {
		t.Error("PageNav should default to false")
	}
	if cfg.InstantNav == nil || *cfg.InstantNav {
		t.Error("InstantNav should default to false")
	}
	if cfg.InlineAssets == nil || *cfg.InlineAssets {
		t.Error("InlineAssets should default to false")
	}
	if cfg.PWA == nil || *cfg.PWA {
		t.Error("PWA should default to false")
	}
	if cfg.Search == nil || *cfg.Search {
		t.Error("Search should default to false")
	}
	if cfg.AllowBrokenLinks == nil || *cfg.AllowBrokenLinks {
		t.Error("AllowBrokenLinks should default to false")
	}
}

func TestMergeConfigs(t *testing.T) {
	defaults := DefaultFileConfig()
	existing := &FileConfig{
		Output:           "./public",
		Title:            "",
		URL:              "https://example.com",
		Theme:            "blog",
		AccentColor:      "#ff6600",
		Port:             IntPtr(8080),
		TopNav:           BoolPtr(true),
		Breadcrumbs:      nil,
		PageNav:          BoolPtr(true),
		InstantNav:       nil,
		InlineAssets:     BoolPtr(true),
		PWA:              BoolPtr(true),
		Search:           nil,
		AllowBrokenLinks: BoolPtr(true),
	}

	merged := MergeConfigs(defaults, existing)
	if merged.Output != "./public" {
		t.Errorf("Output = %q, want %q", merged.Output, "./public")
	}
	if merged.Title != defaults.Title {
		t.Errorf("Title = %q, want %q", merged.Title, defaults.Title)
	}
	if merged.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", merged.URL, "https://example.com")
	}
	if merged.Theme != "blog" {
		t.Errorf("Theme = %q, want %q", merged.Theme, "blog")
	}
	if merged.AccentColor != "#ff6600" {
		t.Errorf("AccentColor = %q, want %q", merged.AccentColor, "#ff6600")
	}
	if merged.Port == nil || *merged.Port != 8080 {
		t.Errorf("Port = %v, want 8080", merged.Port)
	}
	if merged.TopNav == nil || !*merged.TopNav {
		t.Error("TopNav should be true")
	}
	if merged.Breadcrumbs == nil || *merged.Breadcrumbs != *defaults.Breadcrumbs {
		t.Error("Breadcrumbs should remain default")
	}
	if merged.PageNav == nil || !*merged.PageNav {
		t.Error("PageNav should be true")
	}
	if merged.InstantNav == nil || *merged.InstantNav != *defaults.InstantNav {
		t.Error("InstantNav should remain default")
	}
	if merged.InlineAssets == nil || !*merged.InlineAssets {
		t.Error("InlineAssets should be true")
	}
	if merged.PWA == nil || !*merged.PWA {
		t.Error("PWA should be true")
	}
	if merged.Search == nil || *merged.Search != *defaults.Search {
		t.Error("Search should remain default")
	}
	if merged.AllowBrokenLinks == nil || !*merged.AllowBrokenLinks {
		t.Error("AllowBrokenLinks should be true")
	}
}
