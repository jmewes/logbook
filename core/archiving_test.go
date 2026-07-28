package core

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmewes/logbook/config"
)

func Test_Archive_happy_path(t *testing.T) {
	// Arrange
	logBaseDir := createTempDir()
	archiveBaseDir := createTempDir()
	defer func(path string) {
		_ = os.RemoveAll(logBaseDir)
		_ = os.RemoveAll(archiveBaseDir)
	}(logBaseDir)

	logEntry, err := AddLogEntry(logBaseDir, "Log entry for archive test", time.Now())
	if err != nil {
		t.Fatal(err)
	}

	c := config.Configuration{
		LogDirectory:     logBaseDir,
		ArchiveDirectory: archiveBaseDir,
	}

	// Act
	_, err = Archive(c, logEntry.Directory)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	searchResultForLogBaseDir := Search(logBaseDir, "", epoc, nextCentury)
	if len(searchResultForLogBaseDir) != 0 {
		t.Fatal("Expected empty search result")
	}
	searchResultForArchiveBaseDir := Search(archiveBaseDir, "", epoc, nextCentury)
	if len(searchResultForArchiveBaseDir) != 1 {
		t.Fatal("Expected 1 search result")
	}
}

func Test_Archive_redundant_target_directories(t *testing.T) {
	logBaseDir := createTempDir()
	archiveBaseDir := createTempDir()
	defer func(path string) {
		_ = os.RemoveAll(logBaseDir)
		_ = os.RemoveAll(archiveBaseDir)
	}(logBaseDir)

	c := config.Configuration{
		LogDirectory:     logBaseDir,
		ArchiveDirectory: archiveBaseDir,
	}

	now := time.Now()

	// Given an archived "Just a test" logbook entry
	firstEntry, err := AddLogEntry(logBaseDir, "Just a test", now)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Archive(c, firstEntry.Directory)
	if err != nil {
		t.Fatal(err)
	}

	// And a new logbook entry with title "Just a test" has been created
	secondEntry, err := AddLogEntry(logBaseDir, "Just a test", now)
	if err != nil {
		t.Fatal(err)
	}

	// When the logbook entry gets archived
	secondArchivePath, err := Archive(c, secondEntry.Directory)
	if err != nil {
		t.Fatal(err)
	}

	// Then the archive directory name has the suffix "_2"
	if !strings.HasSuffix(secondArchivePath, "_2") {
		t.Fatalf("Expected archive directory to end with '_2', but got: %s", secondArchivePath)
	}
}

func Test_Archive_path_in_subdirectory(t *testing.T) {
	// Arrange
	logBaseDir := createTempDir()
	archiveBaseDir := createTempDir()
	defer func(path string) {
		_ = os.RemoveAll(logBaseDir)
		_ = os.RemoveAll(archiveBaseDir)
	}(logBaseDir)

	logEntry, err := AddLogEntry(logBaseDir, "Log entry for archive test", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	fileInLog := createFileInSubdirectory(logEntry)

	c := config.Configuration{
		LogDirectory:     logBaseDir,
		ArchiveDirectory: archiveBaseDir,
	}

	// Act
	_, err = Archive(c, fileInLog)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	searchResultForLogBaseDir := Search(logBaseDir, "", epoc, nextCentury)
	if len(searchResultForLogBaseDir) != 0 {
		t.Fatal("Expected empty search result")
	}
	searchResultForArchiveBaseDir := Search(archiveBaseDir, "", epoc, nextCentury)
	if len(searchResultForArchiveBaseDir) != 1 {
		t.Fatal("Expected 1 search result")
	}
}
