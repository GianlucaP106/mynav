package core

import (
	"sync"
)

type Container[T comparable] struct {
	container map[string]*T
	mu        *sync.RWMutex
}

func newContainer[T comparable]() *Container[T] {
	return &Container[T]{
		container: make(map[string]*T),
		mu:        &sync.RWMutex{},
	}
}

func (c *Container[T]) Set(key string, d *T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.container[key] = d
}

func (c *Container[T]) Get(key string) *T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.container[key]
}

func (c *Container[T]) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.container, key)
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
	for _, d := range c.container {
		out = append(out, d)
	}
	return out
}

func (c *Container[T]) Contains(key string) bool {
	_, exists := c.container[key]
	return exists
}
