package cmd

import (
	"fmt"
	"io"
)

// Generate handles the static site generation from input folder to output folder
func Generate(cfg *Config, w io.Writer) error {
	_, err := fmt.Fprintf(w, "Generating site from %s to %s with title %q\n", cfg.InputDir, cfg.OutputDir, cfg.Title)
	if err != nil {
		return err
	}
	// TODO: Implement full generation in later stories
	return nil
}
