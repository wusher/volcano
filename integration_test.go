package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationGenerateFromExample(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--title=Integration Test", "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Verify output contains expected messages
	output := stdout.String()
	if !strings.Contains(output, "Generating site") {
		t.Error("Output should contain 'Generating site'")
	}
	if !strings.Contains(output, "Generated") {
		t.Error("Output should contain 'Generated'")
	}

	// Verify expected files exist
	expectedFiles := []string{
		"index.html",
		"404.html",
		"getting-started/index.html",
		"faq/index.html",
		"guides/index.html",
		"guides/installation/index.html",
		"guides/configuration/index.html",
		"api/index.html",
		"api/endpoints/index.html",
	}

	for _, file := range expectedFiles {
		fullPath := filepath.Join(outputDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", file)
		}
	}
}

func TestIntegrationHTMLContent(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "--title=Content Test", "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read the generated index.html
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	html := string(content)

	// Verify HTML structure
	checks := []struct {
		name    string
		content string
	}{
		{"doctype", "<!DOCTYPE html>"},
		{"title", "<title>"},
		{"site title", "Content Test"},
		{"navigation", "class=\"tree-nav\""},
		{"sidebar", "class=\"sidebar\""},
		{"content", "class=\"content\""},
		{"theme toggle", "class=\"theme-toggle\""},
		{"mobile menu", "class=\"mobile-menu-btn\""},
		{"page content", "Welcome to Volcano"},
	}

	for _, check := range checks {
		if !strings.Contains(html, check.content) {
			t.Errorf("index.html should contain %s (%q)", check.name, check.content)
		}
	}
}

func TestIntegrationSyntaxHighlighting(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read a page that has code blocks
	guidesPath := filepath.Join(outputDir, "guides/index.html")
	content, err := os.ReadFile(guidesPath)
	if err != nil {
		t.Fatalf("Failed to read guides/index.html: %v", err)
	}

	html := string(content)

	// Verify code blocks have syntax highlighting classes
	// Chroma adds class="highlight" or class="chroma"
	if !strings.Contains(html, "highlight") && !strings.Contains(html, "chroma") {
		t.Error("Code blocks should have syntax highlighting")
	}
}

func TestIntegrationNavigation(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read index.html
	indexPath := filepath.Join(outputDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	html := string(content)

	// Verify navigation contains links to pages
	navLinks := []string{
		"Getting Started",
		"Guides",
		"Api", // "api" folder becomes "Api" per CleanLabel
	}

	for _, link := range navLinks {
		if !strings.Contains(html, link) {
			t.Errorf("Navigation should contain link to %q", link)
		}
	}
}

func TestIntegration404Page(t *testing.T) {
	// Skip if example folder doesn't exist
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	// Create a temporary output directory
	outputDir := t.TempDir()

	// Run the generator
	var stdout, stderr bytes.Buffer
	exitCode := Run([]string{"-o", outputDir, "./example"}, &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("Generation failed with exit code %d, stderr: %s", exitCode, stderr.String())
	}

	// Read 404.html
	notFoundPath := filepath.Join(outputDir, "404.html")
	content, err := os.ReadFile(notFoundPath)
	if err != nil {
		t.Fatalf("Failed to read 404.html: %v", err)
	}

	html := string(content)

	// Verify 404 page content
	if !strings.Contains(html, "404") {
		t.Error("404 page should contain '404'")
	}
	if !strings.Contains(html, "Page Not Found") {
		t.Error("404 page should contain 'Page Not Found'")
	}
	if !strings.Contains(html, "Return to home") {
		t.Error("404 page should contain link to home")
	}
}
