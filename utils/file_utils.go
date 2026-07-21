package utils

import (
	"path/filepath"
	"strings"
)

func SimpleFileName(path string) string {
	fileName := filepath.Base(path)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func SimpleDirectoryName(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	return parts[len(parts)-1]
}
