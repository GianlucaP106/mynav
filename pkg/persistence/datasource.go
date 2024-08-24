package persistence

import (
	"errors"
	"mynav/pkg/system"
	"sync"
)

type Datasource[T any] struct {
	data *T
	mu   *sync.RWMutex
	Path string
}

func NewDatasource[T any](path string, defaultValue *T) (*Datasource[T], error) {
	ds := &Datasource[T]{
		Path: path,
		mu:   &sync.RWMutex{},
	}

	ds.LoadData()
	if ds.GetData() == nil {
		err := ds.SaveData(defaultValue)
		if err != nil {
			return nil, err
		}
	}

	return ds, nil
}

func (d *Datasource[T]) LoadData() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	data, err := system.LoadJson[T](d.Path)
	if err != nil {
		return errors.New("could not load data from " + d.Path)
	}

	d.data = data
	return nil
}

func (d *Datasource[T]) SaveData(data *T) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = data
	return system.SaveJson(d.data, d.Path)
}

func (d *Datasource[T]) GetData() *T {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.data
}
