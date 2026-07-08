package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const dateLayout = "2006-01-02"

// Config holds validated CLI settings for querying a cookie log.
type Config struct {
	FilePath string
	Date     time.Time
}

// Parse reads -f and -d from args and returns a validated Config.
func Parse(args []string) (Config, error) {
	fs := flag.NewFlagSet("active-cookie", flag.ContinueOnError)
	filePath := fs.String("f", "", "path to the cookie log file (required)")
	dateStr := fs.String("d", "", "date to query in YYYY-MM-DD format, UTC (required)")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if *filePath == "" {
		return Config{}, fmt.Errorf("-f is required")
	}
	if *dateStr == "" {
		return Config{}, fmt.Errorf("-d is required")
	}

	date, err := time.Parse(dateLayout, *dateStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid -d value %q, expected YYYY-MM-DD: %w", *dateStr, err)
	}

	return Config{
		FilePath: *filePath,
		Date:     date,
	}, nil
}

// OpenFile opens the configured cookie log for reading.
func (c Config) OpenFile() (*os.File, error) {
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open %q: %w", c.FilePath, err)
	}
	return file, nil
}
