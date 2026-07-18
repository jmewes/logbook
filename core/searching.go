package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/experimental-software/logbook2/logging"
	"github.com/experimental-software/logbook2/utils"
)

// e.g. /path/to/2026/01/11/18.03_wip
var legacySearchPathPattern = regexp.MustCompile(`.*[/\\](\d{4})[/\\](\d{2})[/\\](\d{2})[/\\](\d{2})\.(\d{2})_.*`)

// e.g. 20250111T1810
var isoDateTimeBasicFormat = regexp.MustCompile(`^\d{4}\d{2}\d{2}T\d{2}\d{2}$`)

func Search(baseDirectory, searchTerm string, from time.Time, to time.Time) []LogbookEntry {
	var result = make([]LogbookEntry, 0)

	maxDepth := strings.Count(baseDirectory, string(os.PathSeparator)) + 4
	normalizedSearchTerm := strings.ToLower(searchTerm)
	searchVisitor := newSearchVisitor(&result, maxDepth, normalizedSearchTerm, from, to)

	if err := filepath.WalkDir(baseDirectory, searchVisitor.collectLogbookEntries); err != nil {
		logging.Error("Could not traverse log directory.", err)
	}
	return result
}

type searchVisitor struct {
	result     *[]LogbookEntry
	maxDepth   int
	searchTerm string
	from       time.Time
	to         time.Time
}

func newSearchVisitor(result *[]LogbookEntry, maxDepth int, searchTerm string, from, to time.Time) searchVisitor {
	return searchVisitor{
		result:     result,
		maxDepth:   maxDepth,
		searchTerm: searchTerm,
		from:       from,
		to:         to,
	}
}

func (v searchVisitor) collectLogbookEntries(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if d.IsDir() && strings.Count(path, string(os.PathSeparator)) > v.maxDepth {
		return fs.SkipDir
	}
	if !isLogEntryFile(path) {
		return nil
	}
	logDatetime := parseLogDatetime(path)

	if !isInRequestedTimeRange(logDatetime, v.from, v.to) {
		return nil
	}

	logEntryFileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	logFile := string(logEntryFileBytes)
	title := strings.Replace(strings.Split(logFile, "\n")[0], "# ", "", 1)
	if v.searchTerm != "" && !strings.Contains(strings.ToLower(title), strings.ToLower(v.searchTerm)) {
		return nil
	}

	logDirectory, _ := filepath.Abs(filepath.Dir(path))
	*v.result = append(*v.result, LogbookEntry{
		DateTime:  logDatetime,
		Directory: logDirectory,
		Title:     title,
	})

	return nil
}

func parseLogDatetime(path string) string {
	if isLegacyTimeNameConvention(path) {
		pathDatetimeMatch := legacySearchPathPattern.FindStringSubmatch(path)
		return fmt.Sprintf("%s-%s-%sT%s:%s",
			pathDatetimeMatch[1],
			pathDatetimeMatch[2],
			pathDatetimeMatch[3],
			pathDatetimeMatch[4],
			pathDatetimeMatch[5],
		)
	}
	simpleFileName := utils.SimpleFileName(path)
	return utils.ToExtendedFormat(simpleFileName)
}

func isInRequestedTimeRange(datetime string, from time.Time, to time.Time) bool {
	return datetime >= from.Format(time.RFC3339) && datetime <= to.Format(time.RFC3339)
}

func isLogEntryFile(path string) bool {
	if !strings.HasSuffix(path, ".md") {
		return false
	}

	pathParts := strings.Split(path, string(os.PathSeparator))
	lastPathPart := pathParts[len(pathParts)-1]
	parentDirectory := pathParts[len(pathParts)-2]

	if regexp.MustCompile(`^\d{2}\.\d{2}_.*`).MatchString(parentDirectory) {
		// keep support for legacy convention where a log entry file has the slug of parent directory suffix, e.g., /path/to/2025/01/14/19.27_foo/foo.md
		return isTimestampBasedLogEntryFile(path)
	}

	simpleFileName := strings.Replace(lastPathPart, ".md", "", -1)
	return len(isoDateTimeBasicFormat.FindStringSubmatch(simpleFileName)) == 1
}

func isTimestampBasedLogEntryFile(path string) bool {
	pathParts := strings.Split(path, string(os.PathSeparator))
	lastPathPart := pathParts[len(pathParts)-1]
	parentDirectory := pathParts[len(pathParts)-2]

	parentDirectorySlug := parentDirectory[6:]
	simpleFileName := strings.Replace(lastPathPart, ".md", "", -1)
	if len(isoDateTimeBasicFormat.FindStringSubmatch(simpleFileName)) != 1 && parentDirectorySlug != simpleFileName {
		return false
	}
	pathDatetimeMatch := legacySearchPathPattern.FindStringSubmatch(path)
	if len(pathDatetimeMatch) == 6 {
		return true
	}

	return false
}

// the legacy convention was to have the log entry time creation as prefix of the parent directory, e.g., /path/to/2025/01/14/19.27_foo
func isLegacyTimeNameConvention(path string) bool {
	pathDatetimeMatch := legacySearchPathPattern.FindStringSubmatch(path)
	if len(pathDatetimeMatch) == 6 {
		return true
	}
	return false
}
