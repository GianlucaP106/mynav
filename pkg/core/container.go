package core

import "sync"

type Container[T any] struct {
	container map[string]*T
	mu        *sync.RWMutex
}

func NewContainer[T any]() *Container[T] {
	return &Container[T]{
		container: make(map[string]*T),
		mu:        &sync.RWMutex{},
	}
}

func (c *Container[T]) Get(id string) *T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.container[id]
}

func (c *Container[T]) Set(key string, d *T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.container[key] = d
}

func (c *Container[T]) SetAll(items []*T, getKey func(*T) string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range items {
		key := getKey(v)
		c.container[key] = v
	}
}

func (c *Container[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.container, key)
}

func (c *Container[T]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.container)
}

func (c Container[T]) All() []*T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]*T, 0)
	for _, d := range c.container {
		out = append(out, d)
	}
	return out
}

type Value[T any] struct {
	val T
	mu  *sync.RWMutex
}

func NewValue[T any](val T) *Value[T] {
	return &Value[T]{
		mu:  &sync.RWMutex{},
		val: val,
	}
}

func (v *Value[T]) Get() T {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.val
}

func (v *Value[T]) Set(val T) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.val = val
}
