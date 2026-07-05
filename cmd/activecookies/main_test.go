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

func writeTempLog(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "cookie_log.csv")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func TestRunSuccess(t *testing.T) {
	path := writeTempLog(t, integrationLog)
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
	path := writeTempLog(t, log)
	var stdout, stderr bytes.Buffer

	err := run([]string{"-f", path, "-d", "2018-12-09"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if !strings.Contains(stderr.String(), "skipping malformed line") {
		t.Errorf("stderr = %q, want malformed line warning", stderr.String())
	}
}
