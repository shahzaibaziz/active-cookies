package app

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

// ParseError describes a single line that could not be parsed, along with
// its position in the file, so a caller can log something actionable.
type ParseError struct {
	Line   int
	Reason error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("line %d: %v", e.Line, e.Reason)
}

// sameUTCDate reports whether two timestamps fall on the same calendar day
// in UTC, which is the only timezone this tool is asked to support.
func sameUTCDate(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}

// FindMostActive reads a cookie log from r and returns every cookie seen the most times
func FindMostActive(r io.Reader, target time.Time) ([]string, []ParseError) {
	c := newCounter()
	var parseErrs []ParseError

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if lineNum == 1 && IsHeaderLine(line) {
			continue
		}

		record, err := ParseLine(line)
		if err != nil {
			parseErrs = append(parseErrs, ParseError{Line: lineNum, Reason: err})
			continue
		}

		if record.Timestamp.UTC().After(target) && !sameUTCDate(record.Timestamp, target) {
			// Still ahead of the target date (log is sorted most-recent-first),
			// keep scanning forward.
			continue
		}

		if sameUTCDate(record.Timestamp, target) {
			c.add(record.Cookie)
			continue
		}

		// Stop: remaining logs are older than the target date.
		break
	}

	return c.mostActive(), parseErrs
}
