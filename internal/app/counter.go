package app

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

func (c *counter) mostActive() []string {
	maxValue := 0
	for _, n := range c.counts {
		if n > maxValue {
			maxValue = n
		}
	}

	if maxValue == 0 {
		return nil
	}

	var result []string
	for _, cookie := range c.order {
		if c.counts[cookie] == maxValue {
			result = append(result, cookie)
		}
	}
	return result
}
