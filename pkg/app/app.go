package app

import (
	"log"
	"mynav/pkg/api"
	"mynav/pkg/ui"
)

type App struct {
	api *api.Api
}

func NewApp() *App {
	api, err := api.NewApi()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &App{
		api: api,
	}
}

func (app *App) Start() {
	ui.Start(app.api)
}
