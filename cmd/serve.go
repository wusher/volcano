package cmd

import (
	"io"

	"volcano/internal/server"
)

// Serve starts the HTTP server to preview the generated site
func Serve(cfg *Config, w io.Writer) error {
	srvConfig := server.Config{
		Dir:     cfg.InputDir,
		Port:    cfg.Port,
		Quiet:   cfg.Quiet,
		Verbose: cfg.Verbose,
	}

	srv := server.New(srvConfig, w)
	return srv.Start()
}
