package app

import "testing"

func TestCounterMostActiveSingleWinner(t *testing.T) {
	c := newCounter()
	c.add("a")
	c.add("b")
	c.add("a")

	got := c.mostActive()
	want := []string{"a"}
	if len(got) != len(want) || got[0] != want[0] {
		t.Errorf("mostActive() = %v, want %v", got, want)
	}
}

func TestCounterMostActiveTiePreservesFirstSeenOrder(t *testing.T) {
	c := newCounter()
	c.add("b")
	c.add("a")
	c.add("b")
	c.add("a")

	got := c.mostActive()
	want := []string{"b", "a"}
	if len(got) != len(want) {
		t.Fatalf("mostActive() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("mostActive()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCounterMostActiveEmpty(t *testing.T) {
	c := newCounter()
	if got := c.mostActive(); got != nil {
		t.Errorf("mostActive() = %v, want nil", got)
	}
}
