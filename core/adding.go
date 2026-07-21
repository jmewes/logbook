package core

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmewes/logbook/utils"
)

var epoc, _ = time.Parse("2006-01-02", "1970-01-01")
var nextCentury, _ = time.Parse("2006-01-02", "2100-01-01")

func AddLogEntry(baseDirectory, title string, dateTime time.Time) (LogbookEntry, error) {
	slug := slugify(title)
	formattedDateTime := fmt.Sprintf("%d-0%d-0%dT0%d:0%d",
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		dateTime.Hour(),
		dateTime.Minute(),
	)

	logDirectoryPath := filepath.Join(baseDirectory,
		fmt.Sprintf("%d", dateTime.Year()),
		fmt.Sprintf("%02d", dateTime.Month()),
		fmt.Sprintf("%02d", dateTime.Day()),
		slug,
	)

	existingEntries := countExistingEntries(logDirectoryPath)
	if existingEntries > 0 {
		logDirectoryPath += "_" + strconv.Itoa(existingEntries+1)
	}

	err := os.MkdirAll(logDirectoryPath, 0777)
	if err != nil {
		return LogbookEntry{}, err
	}

	dateTimeIso8601BasicFormat := fmt.Sprintf("%d%02d%02dT%02d%02d",
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		dateTime.Hour(),
		dateTime.Minute(),
	)
	logEntryFilePath := filepath.Join(logDirectoryPath, dateTimeIso8601BasicFormat+".md")
	err = os.WriteFile(logEntryFilePath, []byte(fmt.Sprintf("# %s\n\n", title)), 0777)
	if err != nil {
		return LogbookEntry{}, err
	}

	return LogbookEntry{DateTime: formattedDateTime, Title: title, Directory: logDirectoryPath}, nil
}

func countExistingEntries(logDirectoryPath string) int {
	parentDirectory := filepath.Dir(logDirectoryPath)
	simpleDirName := utils.SimpleDirectoryName(logDirectoryPath)

	parentDirectoryHandle, err := os.ReadDir(parentDirectory)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range parentDirectoryHandle {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), simpleDirName) {
			count++
		}
	}
	return count
}

func slugify(s string) string {
	result := strings.ToLower(s)

	result = strings.ReplaceAll(result, "ß", "ss")
	result = strings.ReplaceAll(result, "ü", "ue")
	result = strings.ReplaceAll(result, "ö", "oe")
	result = strings.ReplaceAll(result, "ä", "ae")

	result = regexp.MustCompile(`\|`).ReplaceAllString(result, "_")
	result = regexp.MustCompile(`[^A-Za-z0-9_]`).ReplaceAllString(result, "-")
	result = regexp.MustCompile(`-_-`).ReplaceAllString(result, "_")
	result = regexp.MustCompile(`-+`).ReplaceAllString(result, "-")
	result = regexp.MustCompile(`^-`).ReplaceAllString(result, "")
	result = regexp.MustCompile(`-$`).ReplaceAllString(result, "")

	if len(result) > 35 {
		result = result[:35]
	}

	return result
}
