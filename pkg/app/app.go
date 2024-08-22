package app

import (
	"log"
	"mynav/pkg/core"
	"mynav/pkg/logger"
	"mynav/pkg/tasks"
	"mynav/pkg/ui"
)

type App struct {
	api *core.Api
}

func NewApp() *App {
	logger.Init("debug.log")
	api, err := core.NewApi()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &App{
		api: api,
	}
}

func (app *App) Start() {
	tasks.StartExecutor()
	ui.Start(app.api)
}
