package persistence

import (
	"mynav/pkg/system"
	"sync"
)

type Datasource[T any] struct {
	data *T
	mu   *sync.RWMutex
	Path string
}

func NewDatasource[T any](path string) *Datasource[T] {
	ds := &Datasource[T]{
		Path: path,
		mu:   &sync.RWMutex{},
	}

	ds.LoadData()
	return ds
}

func (d *Datasource[T]) LoadData() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = system.LoadJson[T](d.Path)
}

func (d *Datasource[T]) SaveData(data *T) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = data
	system.SaveJson(d.data, d.Path)
}

func (d *Datasource[T]) GetData() *T {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.data
}
