// Package server provides HTTP file server and dynamic rendering functionality.
package server

import (
	"io/fs"
	"os"

	"github.com/wusher/volcano/internal/tree"
)

// Package-level variables for dependency injection (can be overridden in tests)
var (
	osStat     = os.Stat
	osReadFile = os.ReadFile
)

// FileSystem abstracts file system operations for testability
type FileSystem interface {
	// Stat returns file info for the given path
	Stat(path string) (fs.FileInfo, error)
	// ReadFile reads the entire file at path
	ReadFile(path string) ([]byte, error)
}

// TreeScanner abstracts directory tree scanning for testability
type TreeScanner interface {
	// Scan scans a directory and returns the site structure
	Scan(dir string) (*tree.Site, error)
}

// osFileSystem is the default FileSystem using the os package
type osFileSystem struct{}

// Stat implements FileSystem
func (osFileSystem) Stat(path string) (fs.FileInfo, error) {
	return osStat(path)
}

// ReadFile implements FileSystem
func (osFileSystem) ReadFile(path string) ([]byte, error) {
	return osReadFile(path)
}

// defaultScanner is the default TreeScanner using tree.Scan
type defaultScanner struct{}

// Scan implements TreeScanner
func (defaultScanner) Scan(dir string) (*tree.Site, error) {
	return tree.Scan(dir)
}
