package persistence

import "sync"

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
