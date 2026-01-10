package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestIntegrationCLI_FlagOrdering(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Test different flag ordering combinations
	exitCode := Run([]string{"./example", "-o", "/tmp/cli_test"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Failed with flags in different order: %s", stderr.String())
	}

	// Should generate successfully
	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "Generated") {
		t.Error("Should generate site with reordered flags")
	}
}

func TestIntegrationCLI_EqualsSyntax(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Test flag with equals syntax
	exitCode := Run([]string{"-o=/tmp/equals_test", "--title=Test", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Failed with equals syntax: %s", stderr.String())
	}

	// Should generate successfully
	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "Generated") {
		t.Error("Should generate site with equals syntax")
	}
}

func TestIntegrationCLI_VersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Test version flag
	exitCode := Run([]string{"-v"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Failed to show version: %s", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "volcano") && !strings.Contains(output, "version") {
		t.Error("Version output should contain volcano and version")
	}
}

func TestIntegrationCLI_HelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Test help flag
	exitCode := Run([]string{"-h"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Failed to show help: %s", stderr.String())
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "volcano") && !strings.Contains(output, "Usage") {
		t.Error("Help output should contain usage information")
	}
}

func TestIntegrationCLI_QuietFlag(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	var stdout, stderr bytes.Buffer

	// Test quiet flag
	exitCode := Run([]string{"-q", "-o", "/tmp/quiet_test", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Failed with quiet flag: %s", stderr.String())
	}

	// With quiet flag, output should be minimal
	verboseOutput := stdout.String()
	if strings.Contains(verboseOutput, "Generating site") {
		t.Error("Quiet flag should suppress verbose output")
	}
}

func TestIntegrationCLI_VerboseFlag(t *testing.T) {
	if _, err := os.Stat("./example"); os.IsNotExist(err) {
		t.Skip("example folder not found")
	}

	var stdout, stderr bytes.Buffer

	// Test verbose flag
	exitCode := Run([]string{"--verbose", "-o", "/tmp/verbose_test", "./example"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("Failed with verbose flag: %s", stderr.String())
	}

	// With verbose flag, should get detailed output
	verboseOutput := stdout.String()
	if !strings.Contains(verboseOutput, "Generating site") {
		t.Error("Verbose flag should show detailed output")
	}
}

func TestIntegrationCLI_CSSCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Test CSS subcommand
	exitCode := Run([]string{"css"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Errorf("CSS subcommand failed: %s", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "body") || !strings.Contains(output, "color") {
		t.Error("CSS command should output default CSS")
	}
}
