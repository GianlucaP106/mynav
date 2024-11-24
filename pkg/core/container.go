package core

import (
	"sync"
)

type Container[T comparable] struct {
	container map[*T]struct{}
	mu        *sync.RWMutex
}

func newContainer[T comparable]() *Container[T] {
	return &Container[T]{
		container: make(map[*T]struct{}),
		mu:        &sync.RWMutex{},
	}
}

func (c *Container[T]) Add(d *T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.container[d] = struct{}{}
}

func (c *Container[T]) Remove(d *T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.container, d)
}

func (c *Container[T]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.container)
}

func (c *Container[T]) All() []*T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]*T, 0)
	for d := range c.container {
		out = append(out, d)
	}
	return out
}

func (c *Container[T]) Contains(d *T) bool {
	_, exists := c.container[d]
	return exists
}
