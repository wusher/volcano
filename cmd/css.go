package cmd

import (
	"flag"
	"io"
	"os"

	"github.com/wusher/volcano/internal/styles"
)

// CSS outputs the vanilla CSS skeleton to stdout or a file
// Outputs only the theme skeleton (not the layout CSS) for customization
func CSS(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("css", flag.ContinueOnError)
	fs.SetOutput(w)

	var outputFile string
	fs.StringVar(&outputFile, "o", "", "Output file path")
	fs.StringVar(&outputFile, "output", "", "Output file path")

	fs.Usage = func() {
		_, _ = io.WriteString(w, "Usage: volcano css [-o file]\n\n")
		_, _ = io.WriteString(w, "Output the vanilla CSS skeleton for customization.\n\n")
		_, _ = io.WriteString(w, "The output contains only the theme skeleton with helpful comments.\n")
		_, _ = io.WriteString(w, "Layout CSS is automatically included when building your site.\n\n")
		_, _ = io.WriteString(w, "Flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Output only the theme skeleton, not the combined layout + theme
	css := styles.GetThemeCSS("vanilla")

	if outputFile != "" {
		return os.WriteFile(outputFile, []byte(css), 0644)
	}

	_, err := w.Write([]byte(css))
	return err
}
