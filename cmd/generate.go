package cmd

import (
	"io"

	"volcano/internal/generator"
)

// Generate handles the static site generation from input folder to output folder
func Generate(cfg *Config, w io.Writer) error {
	genConfig := generator.Config{
		InputDir:    cfg.InputDir,
		OutputDir:   cfg.OutputDir,
		Title:       cfg.Title,
		Quiet:       cfg.Quiet,
		Verbose:     cfg.Verbose,
		Colored:     cfg.Colored,
		SiteURL:     cfg.SiteURL,
		Author:      cfg.Author,
		OGImage:     cfg.OGImage,
		FaviconPath: cfg.FaviconPath,
		ShowLastMod: cfg.ShowLastMod,
		TopNav:      cfg.TopNav,
	}

	gen, err := generator.New(genConfig, w)
	if err != nil {
		return err
	}

	_, err = gen.Generate()
	return err
}
