package app

// counter tallies cookie occurrences while preserving first-seen order.
// Order is preserved purely so ties resolve deterministically; the spec
// does not require a particular order for tied results, but deterministic
// output is easier to test and reason about than map iteration order.
type counter struct {
	counts map[string]int
	order  []string
}

func newCounter() *counter {
	return &counter{counts: make(map[string]int)}
}

func (c *counter) add(cookie string) {
	if _, seen := c.counts[cookie]; !seen {
		c.order = append(c.order, cookie)
	}
	c.counts[cookie]++
}

// mostActive returns every cookie whose count equals the highest count seen.
// If nothing was added, it returns an empty (nil) slice.
func (c *counter) mostActive() []string {
	max := 0
	for _, n := range c.counts {
		if n > max {
			max = n
		}
	}

	if max == 0 {
		return nil
	}

	var result []string
	for _, cookie := range c.order {
		if c.counts[cookie] == max {
			result = append(result, cookie)
		}
	}
	return result
}
