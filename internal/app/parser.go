package app

import (
	"fmt"
	"strings"
	"time"
)

// Record represents a single cookie sighting in the log.
type Record struct {
	Cookie    string
	Timestamp time.Time
}

const headerLine = "cookie,timestamp"

// ParseLine parses a single CSV line of the form "cookie,timestamp" into a Record
func ParseLine(line string) (Record, error) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return Record{}, fmt.Errorf("empty line")
	}

	fields := strings.Split(trimmed, ",")
	if len(fields) != 2 {
		return Record{}, fmt.Errorf("expected 2 fields (cookie,timestamp), got %d: %q", len(fields), line)
	}

	cookie := strings.TrimSpace(fields[0])
	if cookie == "" {
		return Record{}, fmt.Errorf("empty cookie field: %q", line)
	}

	rawTimestamp := strings.TrimSpace(fields[1])
	timestamp, err := time.Parse(time.RFC3339, rawTimestamp)
	if err != nil {
		return Record{}, fmt.Errorf("invalid timestamp %q: %w", rawTimestamp, err)
	}

	return Record{Cookie: cookie, Timestamp: timestamp}, nil
}

// IsHeaderLine reports whether the given line is the expected CSV header,
func IsHeaderLine(line string) bool {
	return strings.TrimSpace(line) == headerLine
}
