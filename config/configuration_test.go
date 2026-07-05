package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseValidArgs(t *testing.T) {
	cfg, err := Parse([]string{"-f", "cookie_log.csv", "-d", "2018-12-09"})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if cfg.FilePath != "cookie_log.csv" {
		t.Errorf("FilePath = %q, want %q", cfg.FilePath, "cookie_log.csv")
	}

	wantDate := time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)
	if !cfg.Date.Equal(wantDate) {
		t.Errorf("Date = %v, want %v", cfg.Date, wantDate)
	}
}

func TestParseMissingFileFlag(t *testing.T) {
	_, err := Parse([]string{"-d", "2018-12-09"})
	if err == nil {
		t.Fatal("Parse() expected error, got nil")
	}
	if err.Error() != "-f is required" {
		t.Errorf("error = %q, want %q", err.Error(), "-f is required")
	}
}

func TestParseMissingDateFlag(t *testing.T) {
	_, err := Parse([]string{"-f", "cookie_log.csv"})
	if err == nil {
		t.Fatal("Parse() expected error, got nil")
	}
	if err.Error() != "-d is required" {
		t.Errorf("error = %q, want %q", err.Error(), "-d is required")
	}
}

func TestParseInvalidDate(t *testing.T) {
	_, err := Parse([]string{"-f", "cookie_log.csv", "-d", "not-a-date"})
	if err == nil {
		t.Fatal("Parse() expected error, got nil")
	}
	if got, want := err.Error(), `invalid -d value "not-a-date", expected YYYY-MM-DD`; got[:len(want)] != want {
		t.Errorf("error = %q, want prefix %q", err.Error(), want)
	}
}

func TestParseUnknownFlag(t *testing.T) {
	_, err := Parse([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("Parse() expected error, got nil")
	}
}

func TestOpenFileSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.csv")
	if err := os.WriteFile(path, []byte("cookie,timestamp\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg := Config{FilePath: path, Date: time.Now()}
	file, err := cfg.OpenFile()
	if err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}
	file.Close()
}

func TestOpenFileNotFound(t *testing.T) {
	cfg := Config{FilePath: "/no/such/file.csv", Date: time.Now()}
	_, err := cfg.OpenFile()
	if err == nil {
		t.Fatal("OpenFile() expected error, got nil")
	}
	if got, want := err.Error(), `cannot open "/no/such/file.csv"`; got[:len(want)] != want {
		t.Errorf("error = %q, want prefix %q", err.Error(), want)
	}
}
