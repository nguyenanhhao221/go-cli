package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

// filterOut checks whether a file should be filtered out based on its size,
// extension, or whether it is a directory.
//
// Parameters:
//   - path: The full path of the file.
//   - ext: The expected file extension (e.g., ".txt"). If this is an empty string, no extension filtering is applied.
//   - minSize: The minimum file size in bytes. Files smaller than this size will be filtered out.
//   - info: An os.FileInfo object containing file information.
//
// Returns:
//   - true if the file should be filtered out (i.e., ignored); false otherwise.
func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}
	return false
}

func delFile(path string, delLogger *log.Logger) error {
	delLogger.Println(path)
	return os.Remove(path)
}
