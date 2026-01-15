package search

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GenerateSearchIndex writes the search-index.json file to the output directory.
func GenerateSearchIndex(outputDir string, index *Index) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, "search-index.json"), data, 0644)
}
