package ui

import "mynav/pkg/api"

var _api *api.Api

func Api() *api.Api {
	return _api
}

func InitApi() error {
	a, err := api.NewApi()
	if err != nil {
		return err
	}
	_api = a
	return nil
}
