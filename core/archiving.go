package core

import (
	"errors"
	"fmt"
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
	targetDirectoryPath := buildTargetDirectoryPath(configuration, sourceDirectoryPath)

	err = os.MkdirAll(targetDirectoryPath, 0777)
	if err != nil {
		return "", err
	}

	err = gorecurcopy.CopyDirectory(sourceDirectoryPath, targetDirectoryPath)
	if err != nil {
		return "", err
	}

	err = checkAllFilesCopied(sourceDirectoryPath, targetDirectoryPath)
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(sourceDirectoryPath)
	return targetDirectoryPath, err
}

func buildTargetDirectoryPath(configuration config.Configuration, sourceDirectoryPath string) string {
	result := strings.Replace(
		sourceDirectoryPath, configuration.LogDirectory, configuration.ArchiveDirectory, 1,
	)
	result = strings.TrimSuffix(result, "/")

	parentDir := filepath.Dir(result)
	baseName := filepath.Base(result)
	count := 0

	if entries, err := os.ReadDir(parentDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), baseName) {
				count++
			}
		}
	}

	if count > 0 {
		result = fmt.Sprintf("%s_%d", result, count+1)
	}
	return result
}

func checkAllFilesCopied(sourceDirectoryPath string, archiveDirectoryPath string) error {
	sourceFiles := make([]string, 0)
	err := filepath.Walk(sourceDirectoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(sourceDirectoryPath, path)
			if err != nil {
				return err
			}
			sourceFiles = append(sourceFiles, relPath)

		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, sourceFile := range sourceFiles {
		targetPath := filepath.Join(archiveDirectoryPath, sourceFile)
		if _, err := os.Stat(targetPath); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file not copied: %s", sourceFile)
		} else if err != nil {
			return err
		}
	}
	return nil
}
