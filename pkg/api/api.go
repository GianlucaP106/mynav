package api

import (
	"errors"
	"mynav/pkg/core"
	"mynav/pkg/github"
	"mynav/pkg/system"
	"mynav/pkg/tmux"
	"os"
)

type Api struct {
	Tmux   *tmux.TmuxController
	Core   *Core
	Github *github.GithubController
	Port   *system.PortController
}

type Core struct {
	*core.TopicController
	*core.WorkspaceController
	*core.LocalConfiguration
	*core.GlobalConfiguration
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Core = &Core{}
	api.Core.LocalConfiguration = core.NewLocalConfiguration()
	api.Core.GlobalConfiguration = core.NewGlobalConfiguration()

	if api.Core.IsConfigInitialized {
		api.InitControllers()
	} else {
		api.InitTmuxController()
	}

	return api, nil
}

func (api *Api) GetSystemStats() (numTopics int, numWorkspaces int) {
	numTopics = api.Core.GetTopicCount()
	numWorkspaces = api.Core.GetWorkspaceCount()
	return
}

func (api *Api) InitConfiguration() error {
	dir, _ := os.Getwd()
	if _, err := api.Core.InitConfig(dir); err != nil {
		return errors.New("cannot initialize mynav in the home directory")
	}

	api.InitControllers()
	return nil
}

func (api *Api) InitTmuxController() {
	api.Port = system.NewPortController()
	api.Tmux = tmux.NewTmuxController(api.Port)
}

func (api *Api) InitControllers() {
	if api.Core.IsConfigInitialized {
		api.InitTmuxController()
		api.Core.TopicController = core.NewTopicController(api.Core.GetLocalConfigDir(), api.Tmux)
		api.Core.WorkspaceController = core.NewWorkspaceController(
			api.Core.GetTopics(),
			api.Core.GetWorkspaceStorePath(),
			api.Tmux,
			api.Port,
		)

		api.Github = github.NewGithubController(api.Core.GetGithubToken(), func(gat *github.GithubAuthenticationToken) {
			api.Core.SetGithubToken(gat)
		}, func() {
			api.Core.SetGithubToken(nil)
		})
		api.Core.TopicController.WorkspaceController = api.Core.WorkspaceController
	}
}
