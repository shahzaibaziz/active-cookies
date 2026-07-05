package app

import (
	"testing"
	"time"
)

func TestParseLineValid(t *testing.T) {
	record, err := ParseLine("AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00")
	if err != nil {
		t.Fatalf("ParseLine() error = %v", err)
	}

	if record.Cookie != "AtY0laUfhglK3lC7" {
		t.Errorf("Cookie = %q, want %q", record.Cookie, "AtY0laUfhglK3lC7")
	}

	want := time.Date(2018, 12, 9, 14, 19, 0, 0, time.UTC)
	if !record.Timestamp.Equal(want) {
		t.Errorf("Timestamp = %v, want %v", record.Timestamp, want)
	}
}

func TestParseLineTrimsWhitespace(t *testing.T) {
	record, err := ParseLine("  foo , 2018-12-09T00:00:00+00:00  ")
	if err != nil {
		t.Fatalf("ParseLine() error = %v", err)
	}
	if record.Cookie != "foo" {
		t.Errorf("Cookie = %q, want %q", record.Cookie, "foo")
	}
}

func TestParseLineEmptyLine(t *testing.T) {
	_, err := ParseLine("")
	if err == nil {
		t.Fatal("ParseLine() expected error, got nil")
	}
}

func TestParseLineWrongFieldCount(t *testing.T) {
	_, err := ParseLine("only-one-field")
	if err == nil {
		t.Fatal("ParseLine() expected error, got nil")
	}

	_, err = ParseLine("a,b,c")
	if err == nil {
		t.Fatal("ParseLine() expected error for three fields, got nil")
	}
}

func TestParseLineEmptyCookie(t *testing.T) {
	_, err := ParseLine(",2018-12-09T00:00:00+00:00")
	if err == nil {
		t.Fatal("ParseLine() expected error, got nil")
	}
}

func TestParseLineInvalidTimestamp(t *testing.T) {
	_, err := ParseLine("cookie,not-rfc3339")
	if err == nil {
		t.Fatal("ParseLine() expected error, got nil")
	}
}

func TestIsHeaderLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"cookie,timestamp", true},
		{" cookie,timestamp ", true},
		{"Cookie,timestamp", false},
		{"AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00", false},
	}

	for _, tt := range tests {
		if got := IsHeaderLine(tt.line); got != tt.want {
			t.Errorf("IsHeaderLine(%q) = %v, want %v", tt.line, got, tt.want)
		}
	}
}
