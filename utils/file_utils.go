package utils

import (
	"path/filepath"
	"strings"
)

func SimpleFileName(path string) string {
	fileName := filepath.Base(path)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
