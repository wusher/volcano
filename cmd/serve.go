package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"volcano/internal/generator"
	"volcano/internal/server"
)

// Serve starts the HTTP server to preview the generated site
// If the input directory contains markdown files, it generates the site first
func Serve(cfg *Config, w io.Writer) error {
	serveDir := cfg.InputDir

	// Check if this is a source directory (contains .md files but no index.html)
	if isSourceDirectory(cfg.InputDir) {
		// Generate the site first
		_, _ = fmt.Fprintf(w, "Source directory detected, generating site first...\n\n")

		genConfig := generator.Config{
			InputDir:  cfg.InputDir,
			OutputDir: cfg.OutputDir,
			Title:     cfg.Title,
			Quiet:     cfg.Quiet,
			Verbose:   cfg.Verbose,
			Colored:   cfg.Colored,
		}

		gen, err := generator.New(genConfig, w)
		if err != nil {
			return fmt.Errorf("failed to create generator: %w", err)
		}

		_, err = gen.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate site: %w", err)
		}

		_, _ = fmt.Fprintf(w, "\n")
		serveDir = cfg.OutputDir
	}

	srvConfig := server.Config{
		Dir:     serveDir,
		Port:    cfg.Port,
		Quiet:   cfg.Quiet,
		Verbose: cfg.Verbose,
	}

	srv := server.New(srvConfig, w)
	return srv.Start()
}

// isSourceDirectory checks if a directory is a source directory
// (contains .md files but no index.html)
func isSourceDirectory(dir string) bool {
	// Check if index.html exists
	indexHTML := filepath.Join(dir, "index.html")
	if _, err := os.Stat(indexHTML); err == nil {
		return false // Has index.html, it's generated output
	}

	// Check if any .md files exist
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			return true
		}
	}

	return false
}
