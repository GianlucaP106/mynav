package app

import (
	"log"
	"mynav/pkg/core"
	"mynav/pkg/ui"
)

type App struct {
	api *core.Api
}

func newApp() *App {
	api, err := core.NewApi()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &App{
		api: api,
	}
}

func (app *App) start() {
	newCli().run()
	ui.Start(app.api)
}

func Main() {
	newApp().start()
}
