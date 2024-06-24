package core

import "mynav/pkg/system"

type Datasource[T any] struct {
	Data *T
	Path string
}

func NewDatasource[T any](path string) *Datasource[T] {
	ds := &Datasource[T]{
		Path: path,
	}

	ds.LoadData()

	return ds
}

func (d *Datasource[T]) LoadData() {
	d.Data = system.LoadJson[T](d.Path)
}

func (d *Datasource[T]) SaveData() {
	system.SaveJson(d.Data, d.Path)
}
