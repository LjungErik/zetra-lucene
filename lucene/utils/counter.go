package utils

type Counter struct {
	count int
}

func NewCounter() *Counter {
	return &Counter{count: 0}
}

func (c *Counter) GetNextID() int {
	ret := c.count
	c.count += 1
	return ret
}
