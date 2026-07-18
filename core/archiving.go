package core

import (
	"os"
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

	err = os.MkdirAll(targetDirectoryPath, 0777)
	if err != nil {
		return "", err
	}

	err = gorecurcopy.CopyDirectory(sourceDirectoryPath, targetDirectoryPath)
	if err != nil {
		return "", err
	}
	err = os.RemoveAll(sourceDirectoryPath)
	return targetDirectoryPath, err
}
