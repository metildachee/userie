package elasticsearch

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) GetCount() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
	return fmt.Sprintf("%d", c.count)
}
