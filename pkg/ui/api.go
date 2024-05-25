package ui

import "mynav/pkg/api"

var _api *api.Api

func Api() *api.Api {
	return _api
}

func InitApi() {
	_api = api.NewApi()
}
