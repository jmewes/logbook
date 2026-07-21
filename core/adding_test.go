package core

import (
	"os"
	"testing"
	"time"

	"github.com/jmewes/logbook/utils"
)

func TestAddLogEntry(t *testing.T) {
	// Arrange
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer os.RemoveAll(tempDir)

	// Act
	entry, err := AddLogEntry(tempDir, "This is a new log entry", time.Now())

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if entry.Title != "This is a new log entry" {
		t.Errorf("AddLogEntry returned wrong title")
	}
	allLogEntries := Search(tempDir, "", epoc, nextCentury)
	if len(allLogEntries) != 1 {
		t.Errorf("AddLogEntry returned wrong number of entries")
	}
}

// should add counter-suffix for duplicate titles
func TestAddLogEntry_duplicate_title(t *testing.T) {
	fixture := setup(t)
	t.Cleanup(fixture.cleanup)
	// Given a logbook entry with the title "Foo"
	_, _ = AddLogEntry(fixture.tempDir, "Foo", time.Now())

	// When another logbook entry with the title "Foo" gets created
	entry, _ := AddLogEntry(fixture.tempDir, "Foo", time.Now())

	// Then the directory for the second logbook entry has the suffix "_2"
	dir := utils.SimpleDirectoryName(entry.Directory)
	if dir != "foo_2" {
		t.Errorf("actual: %v, expected: %v", dir, "foo_2")
	}
}

func setup(t *testing.T) fixture {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}

	return fixture{
		cleanup: func() {
			//goland:noinspection GoUnhandledErrorResult
			os.RemoveAll(tempDir)
		},
		tempDir: tempDir,
	}
}

type fixture struct {
	cleanup func()
	tempDir string
}
