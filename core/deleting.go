package core

import (
	"github.com/laurent22/go-trash"
)

func Remove(sourcePath string) error {
	sourceDirectoryPath, err := logbookEntryRootPath(sourcePath)
	if err != nil {
		return err
	}

	_, err = trash.MoveToTrash(sourceDirectoryPath)

	return err
}
