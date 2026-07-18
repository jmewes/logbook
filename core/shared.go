package core

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type LogbookEntry struct {
	DateTime  string `json:"dateTime"`
	Title     string `json:"title"`
	Directory string `json:"directory"`
}

func logbookEntryRootPath(path string) (string, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	re := regexp.MustCompile(`(.*[/\\]\d{4}[/\\]\d{2}[/\\]\d{2}[/\\]\d{2}\.\d{2}_.*?[/\\]).*`)
	m := re.FindStringSubmatch(path)
	if len(m) == 2 {
		return m[1], nil
	}
	re = regexp.MustCompile(`(.*[/\\]\d{4}[/\\]\d{2}[/\\]\d{2}[/\\].*?[/\\]).*`)
	m = re.FindStringSubmatch(path)
	if len(m) == 2 && containsLogbookEntryFile(m[1]) {
		return m[1], nil
	}
	return "", errors.New("invalid logbook entry path: " + path)
}

func containsLogbookEntryFile(dirPath string) bool {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false
	}

	iso8601Pattern := regexp.MustCompile(`^\d{8}T\d{4}\.md$`)
	for _, entry := range entries {
		if !entry.IsDir() && iso8601Pattern.MatchString(entry.Name()) {
			return true
		}
	}
	return false
}
