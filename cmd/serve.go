package cmd

import (
	"fmt"
	"io"
)

// Serve starts the HTTP server to preview the generated site
func Serve(cfg *Config, w io.Writer) error {
	_, err := fmt.Fprintf(w, "Serving %s on port %d\n", cfg.InputDir, cfg.Port)
	if err != nil {
		return err
	}
	// TODO: Implement full server in later stories
	return nil
}
