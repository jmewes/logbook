package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmewes/logbook/config"
	"github.com/plus3it/gorecurcopy"
)

func Archive(configuration config.Configuration, sourcePath string) (string, error) {
	sourceDirectoryPath, err := logbookEntryRootPath(sourcePath)
	if err != nil {
		return "", err
	}
	targetDirectoryPath := strings.Replace(
		sourceDirectoryPath, configuration.LogDirectory, configuration.ArchiveDirectory, 1,
	)

	if _, err := os.Stat(targetDirectoryPath); err == nil {
		return "", errors.New("target directory already exists: " + targetDirectoryPath)
	}

	err = os.MkdirAll(targetDirectoryPath, 0777)
	if err != nil {
		return "", err
	}

	err = gorecurcopy.CopyDirectory(sourceDirectoryPath, targetDirectoryPath)
	if err != nil {
		return "", err
	}

	// Verify that all files have been copied
	sourceFiles := make(map[string]os.FileInfo)
	err = filepath.Walk(sourceDirectoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(sourceDirectoryPath, path)
			if err != nil {
				return err
			}
			sourceFiles[relPath] = info
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	for relPath, sourceInfo := range sourceFiles {
		targetPath := filepath.Join(targetDirectoryPath, relPath)
		targetInfo, err := os.Stat(targetPath)
		if err != nil {
			return "", err
		}
		if sourceInfo.Size() != targetInfo.Size() {
			return "", errors.New("file size mismatch: " + relPath)
		}
	}

	err = os.RemoveAll(sourceDirectoryPath)
	return targetDirectoryPath, err
}
