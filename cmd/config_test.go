package cmd

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OutputDir != "./output" {
		t.Errorf("OutputDir = %q, want %q", cfg.OutputDir, "./output")
	}

	if cfg.Port != 1776 {
		t.Errorf("Port = %d, want %d", cfg.Port, 1776)
	}

	if cfg.Title != "My Site" {
		t.Errorf("Title = %q, want %q", cfg.Title, "My Site")
	}

	if cfg.ServeMode {
		t.Error("ServeMode should be false by default")
	}

	if cfg.Quiet {
		t.Error("Quiet should be false by default")
	}

	if cfg.Verbose {
		t.Error("Verbose should be false by default")
	}

	if !cfg.ShowBreadcrumbs {
		t.Error("ShowBreadcrumbs should be true by default")
	}

	if !cfg.ViewTransitions {
		t.Error("ViewTransitions should be true by default")
	}

	if cfg.Colored {
		t.Error("Colored should be false by default")
	}

	if cfg.AllowBrokenLinks {
		t.Error("AllowBrokenLinks should be false by default")
	}
}
