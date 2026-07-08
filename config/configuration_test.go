package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTestCSV(t *testing.T, name, content string) string {
	t.Helper()

	root := projectRoot(t)
	dir := filepath.Join(root, "tmp")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	safeName := strings.ReplaceAll(t.Name(), "/", "_") + "_" + name
	path := filepath.Join(dir, safeName)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func projectRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}

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
	path := writeTestCSV(t, "log.csv", "cookie,timestamp\n")

	cfg := Config{FilePath: path, Date: time.Now()}
	file, err := cfg.OpenFile()
	if err != nil {
		t.Fatalf("OpenFile() error = %v", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
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
