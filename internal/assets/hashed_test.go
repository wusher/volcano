package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteHashedAsset(t *testing.T) {
	t.Run("basic asset", func(t *testing.T) {
		tmpDir := t.TempDir()
		content := "body { color: red; }"

		asset, err := WriteHashedAsset(tmpDir, "styles", "css", content, "")
		if err != nil {
			t.Fatalf("WriteHashedAsset() error = %v", err)
		}

		if asset == nil {
			t.Fatal("WriteHashedAsset() returned nil")
		}

		// Check filename format
		if !strings.HasPrefix(asset.FileName, "styles.") {
			t.Errorf("FileName should start with 'styles.', got %q", asset.FileName)
		}
		if !strings.HasSuffix(asset.FileName, ".css") {
			t.Errorf("FileName should end with '.css', got %q", asset.FileName)
		}

		// Check URL path
		if !strings.HasPrefix(asset.URLPath, "/assets/styles.") {
			t.Errorf("URLPath should start with '/assets/styles.', got %q", asset.URLPath)
		}

		// Check content
		if asset.Content != content {
			t.Errorf("Content = %q, want %q", asset.Content, content)
		}

		// Verify file was written
		filePath := filepath.Join(tmpDir, "assets", asset.FileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read asset file: %v", err)
		}
		if string(data) != content {
			t.Errorf("File content = %q, want %q", string(data), content)
		}
	})

	t.Run("with base URL", func(t *testing.T) {
		tmpDir := t.TempDir()
		content := "body { color: blue; }"

		asset, err := WriteHashedAsset(tmpDir, "app", "js", content, "/docs")
		if err != nil {
			t.Fatalf("WriteHashedAsset() error = %v", err)
		}

		if !strings.HasPrefix(asset.URLPath, "/docs/assets/") {
			t.Errorf("URLPath should include base URL, got %q", asset.URLPath)
		}
	})

	t.Run("deterministic hash", func(t *testing.T) {
		tmpDir1 := t.TempDir()
		tmpDir2 := t.TempDir()
		content := "identical content"

		asset1, _ := WriteHashedAsset(tmpDir1, "test", "txt", content, "")
		asset2, _ := WriteHashedAsset(tmpDir2, "test", "txt", content, "")

		if asset1.FileName != asset2.FileName {
			t.Errorf("Same content should produce same hash: %q != %q", asset1.FileName, asset2.FileName)
		}
	})

	t.Run("different content produces different hash", func(t *testing.T) {
		tmpDir := t.TempDir()

		asset1, _ := WriteHashedAsset(tmpDir, "test", "txt", "content1", "")
		// Need different name or different dir since same file would overwrite
		tmpDir2 := t.TempDir()
		asset2, _ := WriteHashedAsset(tmpDir2, "test", "txt", "content2", "")

		if asset1.FileName == asset2.FileName {
			t.Error("Different content should produce different hash")
		}
	})

	t.Run("creates assets directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		_, err := WriteHashedAsset(tmpDir, "styles", "css", "body{}", "")
		if err != nil {
			t.Fatalf("WriteHashedAsset() error = %v", err)
		}

		assetsDir := filepath.Join(tmpDir, "assets")
		info, err := os.Stat(assetsDir)
		if err != nil {
			t.Fatalf("Assets directory not created: %v", err)
		}
		if !info.IsDir() {
			t.Error("assets should be a directory")
		}
	})
}
