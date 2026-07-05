package app

import (
	"strings"
	"testing"
	"time"
)

const sampleLog = `cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-09T10:13:00+00:00
5UAVanZf6UtGyKVS,2018-12-09T07:25:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
SAZuXPGUrfbcn5UA,2018-12-08T22:03:00+00:00
4sMM2LxV07bPJzwf,2018-12-08T21:30:00+00:00
fbcn5UAVanZf6UtG,2018-12-08T09:30:00+00:00
4sMM2LxV07bPJzwf,2018-12-07T23:30:00+00:00
`

func TestFindMostActiveSingleWinner(t *testing.T) {
	target := time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)
	cookies, parseErrs := FindMostActive(strings.NewReader(sampleLog), target)

	if len(parseErrs) != 0 {
		t.Fatalf("parse errors = %v, want none", parseErrs)
	}

	want := []string{"AtY0laUfhglK3lC7"}
	if len(cookies) != len(want) || cookies[0] != want[0] {
		t.Errorf("FindMostActive() = %v, want %v", cookies, want)
	}
}

func TestFindMostActiveTie(t *testing.T) {
	target := time.Date(2018, 12, 8, 0, 0, 0, 0, time.UTC)
	cookies, parseErrs := FindMostActive(strings.NewReader(sampleLog), target)

	if len(parseErrs) != 0 {
		t.Fatalf("parse errors = %v, want none", parseErrs)
	}

	want := []string{"SAZuXPGUrfbcn5UA", "4sMM2LxV07bPJzwf", "fbcn5UAVanZf6UtG"}
	if len(cookies) != len(want) {
		t.Fatalf("FindMostActive() = %v, want %v", cookies, want)
	}
	for i := range want {
		if cookies[i] != want[i] {
			t.Errorf("cookies[%d] = %q, want %q", i, cookies[i], want[i])
		}
	}
}

func TestFindMostActiveNoMatches(t *testing.T) {
	target := time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC)
	cookies, parseErrs := FindMostActive(strings.NewReader(sampleLog), target)

	if len(parseErrs) != 0 {
		t.Fatalf("parse errors = %v, want none", parseErrs)
	}
	if cookies != nil {
		t.Errorf("FindMostActive() = %v, want nil", cookies)
	}
}

func TestFindMostActiveSkipsMalformedLines(t *testing.T) {
	log := `cookie,timestamp
bad-line
,2018-12-09T00:00:00+00:00
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
AtY0laUfhglK3lC7,2018-12-08T06:19:00+00:00
`
	target := time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)

	cookies, parseErrs := FindMostActive(strings.NewReader(log), target)
	if len(parseErrs) != 2 {
		t.Fatalf("parse errors = %d, want 2", len(parseErrs))
	}
	if len(cookies) != 1 || cookies[0] != "AtY0laUfhglK3lC7" {
		t.Errorf("FindMostActive() = %v, want [AtY0laUfhglK3lC7]", cookies)
	}
}

func TestFindMostActiveWithoutHeader(t *testing.T) {
	log := `AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
AtY0laUfhglK3lC7,2018-12-09T06:19:00+00:00
`
	target := time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)
	cookies, parseErrs := FindMostActive(strings.NewReader(log), target)

	if len(parseErrs) != 0 {
		t.Fatalf("parse errors = %v, want none", parseErrs)
	}
	if len(cookies) != 1 || cookies[0] != "AtY0laUfhglK3lC7" {
		t.Errorf("FindMostActive() = %v, want [AtY0laUfhglK3lC7]", cookies)
	}
}

func TestFindMostActiveStopsEarly(t *testing.T) {
	log := `AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
AtY0laUfhglK3lC7,2018-12-08T06:19:00+00:00
ignored,2018-12-07T06:19:00+00:00
`
	target := time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC)
	cookies, parseErrs := FindMostActive(strings.NewReader(log), target)

	if len(parseErrs) != 0 {
		t.Fatalf("parse errors = %v, want none", parseErrs)
	}
	if len(cookies) != 1 || cookies[0] != "AtY0laUfhglK3lC7" {
		t.Errorf("FindMostActive() = %v, want [AtY0laUfhglK3lC7]", cookies)
	}
}

func TestSameUTCDate(t *testing.T) {
	a := time.Date(2018, 12, 9, 2, 0, 0, 0, time.FixedZone("+2", 2*3600)) // 2018-12-09 00:00 UTC
	b := time.Date(2018, 12, 9, 1, 0, 0, 0, time.UTC)
	if !sameUTCDate(a, b) {
		t.Error("sameUTCDate() = false, want true for same UTC calendar day")
	}
}
