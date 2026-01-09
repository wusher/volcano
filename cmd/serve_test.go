package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestServe(t *testing.T) {
	cfg := &Config{
		InputDir: "/tmp/test-serve",
		Port:     8080,
	}

	var buf bytes.Buffer
	err := Serve(cfg, &buf)
	if err != nil {
		t.Errorf("Serve() unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Serving") {
		t.Error("Serve should print 'Serving'")
	}
	if !strings.Contains(output, cfg.InputDir) {
		t.Error("Serve should print input directory")
	}
	if !strings.Contains(output, "8080") {
		t.Error("Serve should print port number")
	}
}
