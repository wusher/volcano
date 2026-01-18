package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wusher/volcano/internal/config"
	"github.com/wusher/volcano/internal/output"
)

func TestConfigTrackerSetAndGet(t *testing.T) {
	tracker := newConfigTracker()
	tracker.set("title", "My Site", sourceFile)

	if tracker.values["title"] != "My Site" {
		t.Errorf("expected stored value, got %v", tracker.values["title"])
	}
	if tracker.getSource("title") != sourceFile {
		t.Errorf("expected sourceFile, got %v", tracker.getSource("title"))
	}
}

func TestInitDefaultTracker(t *testing.T) {
	cfg := DefaultConfig()
	tracker := newConfigTracker()
	initDefaultTracker(tracker, cfg)

	if tracker.getSource("output") != sourceDefault {
		t.Errorf("expected output source default, got %v", tracker.getSource("output"))
	}
	if tracker.values["output"] != cfg.OutputDir {
		t.Errorf("expected output value %q, got %v", cfg.OutputDir, tracker.values["output"])
	}
	if tracker.getSource("breadcrumbs") != sourceDefault {
		t.Errorf("expected breadcrumbs source default, got %v", tracker.getSource("breadcrumbs"))
	}
}

func TestDetectCLIOverrides(t *testing.T) {
	cfg := DefaultConfig()
	pre := copyConfigValues(cfg)

	cfg.Title = "New Title"
	cfg.PWA = true

	tracker := newConfigTracker()
	detectCLIOverrides(cfg, pre, tracker)

	if tracker.getSource("title") != sourceCLI {
		t.Errorf("expected title source CLI, got %v", tracker.getSource("title"))
	}
	if tracker.getSource("pwa") != sourceCLI {
		t.Errorf("expected pwa source CLI, got %v", tracker.getSource("pwa"))
	}
}

func TestPrintCLIOverrides(t *testing.T) {
	tracker := newConfigTracker()
	tracker.set("title", "Custom", sourceCLI)

	var buf bytes.Buffer
	logger := output.NewLogger(&buf, false, false, false)
	printCLIOverrides(logger, tracker, true)

	output := buf.String()
	if !strings.Contains(output, "CLI --title overrides config file") {
		t.Errorf("expected override message, got %q", output)
	}

	buf.Reset()
	printCLIOverrides(logger, tracker, false)
	if buf.String() != "" {
		t.Errorf("expected no output without config file, got %q", buf.String())
	}
}

func TestPrintBuildConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.InputDir = "./docs"
	cfg.OutputDir = "./output"
	cfg.SiteURL = "https://example.com"
	cfg.Title = "Docs"
	cfg.Author = "Author"
	cfg.Theme = "docs"
	cfg.CSSPath = "./custom.css"
	cfg.AccentColor = "#ff6600"
	cfg.FaviconPath = "./favicon.ico"
	cfg.OGImage = "./og.png"
	cfg.TopNav = true
	cfg.ShowBreadcrumbs = true
	cfg.ShowPageNav = true
	cfg.InstantNav = true
	cfg.InlineAssets = true
	cfg.PWA = true
	cfg.Search = true
	cfg.AllowBrokenLinks = true

	var buf bytes.Buffer
	logger := output.NewLogger(&buf, false, false, false)
	printBuildConfig(logger, cfg)

	output := buf.String()
	if !strings.Contains(output, "input:") || !strings.Contains(output, "features:") {
		t.Errorf("expected build config output, got %q", output)
	}
	if !strings.Contains(output, "allowBrokenLinks") {
		t.Errorf("expected allowBrokenLinks in features, got %q", output)
	}
}

func TestPrescanServeArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedInput  string
		expectedConfig string
	}{
		{
			name:           "input only",
			args:           []string{"./docs"},
			expectedInput:  "./docs",
			expectedConfig: "",
		},
		{
			name:           "config flag with space",
			args:           []string{"--config", "./config.json", "./docs"},
			expectedInput:  "./docs",
			expectedConfig: "./config.json",
		},
		{
			name:           "config flag with equals",
			args:           []string{"--config=./config.json", "./docs"},
			expectedInput:  "./docs",
			expectedConfig: "./config.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inputDir, configPath := prescanServeArgs(tc.args)
			if inputDir != tc.expectedInput {
				t.Errorf("inputDir = %q, want %q", inputDir, tc.expectedInput)
			}
			if configPath != tc.expectedConfig {
				t.Errorf("configPath = %q, want %q", configPath, tc.expectedConfig)
			}
		})
	}
}

func TestApplyServeFileConfig(t *testing.T) {
	cfg := DefaultConfig()
	tracker := newConfigTracker()
	fileCfg := &config.FileConfig{
		Port:        config.IntPtr(8080),
		Title:       "Serve Title",
		Theme:       "blog",
		TopNav:      config.BoolPtr(true),
		Breadcrumbs: config.BoolPtr(false),
		PWA:         config.BoolPtr(true),
		Search:      config.BoolPtr(true),
	}

	applyServeFileConfig(cfg, fileCfg, tracker)

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want %d", cfg.Port, 8080)
	}
	if cfg.Title != "Serve Title" {
		t.Errorf("Title = %q, want %q", cfg.Title, "Serve Title")
	}
	if cfg.Theme != "blog" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "blog")
	}
	if !cfg.TopNav {
		t.Error("TopNav should be true")
	}
	if cfg.ShowBreadcrumbs {
		t.Error("ShowBreadcrumbs should be false")
	}
	if !cfg.PWA {
		t.Error("PWA should be true")
	}
	if !cfg.Search {
		t.Error("Search should be true")
	}
	if tracker.getSource("title") != sourceFile {
		t.Errorf("expected title source file, got %v", tracker.getSource("title"))
	}
}

func TestInitServeDefaultTracker(t *testing.T) {
	cfg := DefaultConfig()
	tracker := newConfigTracker()
	initServeDefaultTracker(tracker, cfg)

	if tracker.getSource("port") != sourceDefault {
		t.Errorf("expected port source default, got %v", tracker.getSource("port"))
	}
	if tracker.values["port"] != cfg.Port {
		t.Errorf("expected port value %d, got %v", cfg.Port, tracker.values["port"])
	}
}

func TestDetectServeCLIOverrides(t *testing.T) {
	cfg := DefaultConfig()
	pre := copyServeConfigValues(cfg)
	cfg.Port = 8081
	cfg.Search = true

	tracker := newConfigTracker()
	detectServeCLIOverrides(cfg, pre, tracker)

	if tracker.getSource("port") != sourceCLI {
		t.Errorf("expected port source CLI, got %v", tracker.getSource("port"))
	}
	if tracker.getSource("search") != sourceCLI {
		t.Errorf("expected search source CLI, got %v", tracker.getSource("search"))
	}
}

func TestPrintServeCLIOverrides(t *testing.T) {
	tracker := newConfigTracker()
	tracker.set("port", 8080, sourceCLI)

	var buf bytes.Buffer
	logger := output.NewLogger(&buf, false, false, false)
	printServeCLIOverrides(logger, tracker, true)

	output := buf.String()
	if !strings.Contains(output, "CLI --port overrides config file") {
		t.Errorf("expected override message, got %q", output)
	}
}

func TestPrintServeConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.InputDir = "./docs"
	cfg.Port = 8080
	cfg.Title = "Docs"
	cfg.SiteURL = "https://example.com"
	cfg.Author = "Author"
	cfg.Theme = "docs"
	cfg.CSSPath = "./custom.css"
	cfg.AccentColor = "#ff6600"
	cfg.FaviconPath = "./favicon.ico"
	cfg.TopNav = true
	cfg.ShowBreadcrumbs = true
	cfg.ShowPageNav = true
	cfg.InstantNav = true
	cfg.PWA = true
	cfg.Search = true

	var buf bytes.Buffer
	logger := output.NewLogger(&buf, false, false, false)
	printServeConfig(logger, cfg)

	output := buf.String()
	if !strings.Contains(output, "port:") || !strings.Contains(output, "features:") {
		t.Errorf("expected serve config output, got %q", output)
	}
	if !strings.Contains(output, "search") {
		t.Errorf("expected search in features, got %q", output)
	}
}
