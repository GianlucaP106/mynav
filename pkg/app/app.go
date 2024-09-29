package app

import (
	"fmt"
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

	ws := app.api.Workspaces.GetWorkspaces()
	w := ws[0]
	w2, err := app.api.Workspaces.CreateSubworkspace("ss", w)
	fmt.Println(w2, err)

	ui.Start(app.api)
}

func Main() {
	newApp().start()
}
