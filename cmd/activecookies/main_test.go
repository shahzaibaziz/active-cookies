package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const integrationLog = `cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
`

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

func TestRunSuccess(t *testing.T) {
	path := writeTestCSV(t, "cookie_log.csv", integrationLog)
	var stdout, stderr bytes.Buffer

	err := run([]string{"-f", path, "-d", "2018-12-09"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	got := strings.TrimSpace(stdout.String())
	want := "AtY0laUfhglK3lC7"
	if got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunMissingFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run() expected error, got nil")
	}
}

func TestRunFileNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := run([]string{"-f", "/no/such/file.csv", "-d", "2018-12-09"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run() expected error, got nil")
	}
}

func TestRunLogsMalformedLines(t *testing.T) {
	log := integrationLog + "bad-line\n"
	path := writeTestCSV(t, "cookie_log.csv", log)
	var stdout, stderr bytes.Buffer

	err := run([]string{"-f", path, "-d", "2018-12-09"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if !strings.Contains(stderr.String(), "skipping malformed line") {
		t.Errorf("stderr = %q, want malformed line warning", stderr.String())
	}
}
