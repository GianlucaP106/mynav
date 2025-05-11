package core

import (
	"errors"
	"sync"
)

type Datasource[T any] struct {
	data *T
	mu   *sync.RWMutex
	Path string
}

func newDatasource[T any](path string, defaultValue *T) (*Datasource[T], error) {
	ds := &Datasource[T]{
		Path: path,
		mu:   &sync.RWMutex{},
	}

	ds.Load()
	if ds.Get() == nil {
		ds.data = defaultValue
	}

	return ds, nil
}

func (d *Datasource[T]) Load() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	data, err := LoadJson[T](d.Path)
	if err != nil {
		return errors.New("could not load data from " + d.Path)
	}

	d.data = data
	return nil
}

func (d *Datasource[T]) Save(data *T) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = data
	return SaveJson(d.data, d.Path)
}

func (d *Datasource[T]) Get() *T {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.data
}
