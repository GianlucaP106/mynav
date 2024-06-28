package ui

import "mynav/pkg/core"

var _api *core.Api

func Api() *core.Api {
	return _api
}

func InitApi() error {
	a, err := core.NewApi()
	if err != nil {
		return err
	}
	_api = a
	return nil
}
