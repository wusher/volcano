package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wusher/volcano/internal/config"
	"github.com/wusher/volcano/internal/output"
)

// Init handles the init subcommand for creating/updating a config file
func Init(args []string, stdout, stderr io.Writer) error {
	errLogger := output.NewLogger(stderr, output.IsStderrTTY(), false, false)
	stdLogger := output.NewLogger(stdout, output.IsStdoutTTY(), false, false)

	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var showHelp bool
	var outputPath string

	fs.StringVar(&outputPath, "o", "", "Output path for config file (default: ./volcano.json)")
	fs.StringVar(&outputPath, "output", "", "Output path for config file (default: ./volcano.json)")
	fs.BoolVar(&showHelp, "h", false, "Show help")
	fs.BoolVar(&showHelp, "help", false, "Show help")

	fs.Usage = func() {
		printInitUsage(stdout)
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if showHelp {
		printInitUsage(stdout)
		return nil
	}

	// Default output path is current directory
	if outputPath == "" {
		outputPath = config.ConfigFileName
	}

	// Check if the path is a directory, if so append the config filename
	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		outputPath = filepath.Join(outputPath, config.ConfigFileName)
	}

	// Check if file already exists
	var existingConfig *config.FileConfig
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, load it
		existingConfig, err = config.Load(outputPath)
		if err != nil {
			errLogger.Error("failed to read existing config: %v", err)
			return err
		}
		stdLogger.Println("Updating existing config file: %s", outputPath)
	} else {
		stdLogger.Println("Creating new config file: %s", outputPath)
	}

	// Create default config with all options
	newConfig := config.DefaultFileConfig()

	// If existing config, merge it (existing values override defaults)
	if existingConfig != nil {
		newConfig = config.MergeConfigs(newConfig, existingConfig)
	}

	// Write the config file with pretty formatting
	data, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		errLogger.Error("failed to marshal config: %v", err)
		return err
	}

	// Add trailing newline
	data = append(data, '\n')

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		errLogger.Error("failed to write config file: %v", err)
		return err
	}

	stdLogger.Success("Config file written: %s", outputPath)
	return nil
}

func printInitUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Create or update a volcano.json configuration file")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  volcano init [flags] [directory]")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Flags:")
	_, _ = fmt.Fprintln(w, "  -o, --output <path>  Output path (default: ./volcano.json)")
	_, _ = fmt.Fprintln(w, "  -h, --help           Show this help message")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Description:")
	_, _ = fmt.Fprintln(w, "  Creates a new volcano.json with all available options and their defaults.")
	_, _ = fmt.Fprintln(w, "  If a config file already exists, adds any missing keys while preserving")
	_, _ = fmt.Fprintln(w, "  existing values.")
	_, _ = fmt.Fprintln(w, "")
	_, _ = fmt.Fprintln(w, "Examples:")
	_, _ = fmt.Fprintln(w, "  volcano init                      Create volcano.json in current directory")
	_, _ = fmt.Fprintln(w, "  volcano init -o ./docs            Create volcano.json in ./docs directory")
	_, _ = fmt.Fprintln(w, "  volcano init -o my-config.json    Create config with custom name")
}
