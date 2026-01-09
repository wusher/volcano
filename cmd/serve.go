package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"volcano/internal/server"
)

// Serve starts the HTTP server to preview the site
// If the input directory is a source directory (contains .md files), it uses
// dynamic rendering so changes are reflected immediately without restart.
// Otherwise, it serves static files from the directory.
func Serve(cfg *Config, w io.Writer) error {
	// Check if this is a source directory (contains .md files but no index.html)
	if isSourceDirectory(cfg.InputDir) {
		// Use dynamic server for live rendering
		dynamicCfg := server.DynamicConfig{
			SourceDir: cfg.InputDir,
			Title:     cfg.Title,
			Port:      cfg.Port,
			Quiet:     cfg.Quiet,
			Verbose:   cfg.Verbose,
		}

		srv, err := server.NewDynamicServer(dynamicCfg, w)
		if err != nil {
			return err
		}

		return srv.Start()
	}

	// Serve static files from the directory
	srvConfig := server.Config{
		Dir:     cfg.InputDir,
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
