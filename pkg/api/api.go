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
	Tmux          *tmux.TmuxController
	Core          *Core
	Configuration *Configuration
	Github        *github.GithubController
	Port          *system.PortController
}

type Configuration struct {
	*core.LocalConfiguration
	*core.GlobalConfiguration
}

type Core struct {
	*core.TopicController
	*core.WorkspaceController
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Core = &Core{}
	api.Configuration = &Configuration{}
	api.Configuration.LocalConfiguration = core.NewLocalConfiguration()
	api.Configuration.GlobalConfiguration = core.NewGlobalConfiguration()

	if api.Configuration.IsConfigInitialized {
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
	if _, err := api.Configuration.InitConfig(dir); err != nil {
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
	if api.Configuration.IsConfigInitialized {
		api.InitTmuxController()
		api.Core.TopicController = core.NewTopicController(api.Configuration.GetLocalConfigDir(), api.Tmux)
		api.Core.WorkspaceController = core.NewWorkspaceController(
			api.Core.GetTopics(),
			api.Configuration.GetWorkspaceStorePath(),
			api.Tmux,
			api.Port,
		)

		api.Github = github.NewGithubController(api.Configuration.GetGithubToken(), func(gat *github.GithubAuthenticationToken) {
			api.Configuration.SetGithubToken(gat)
		}, func() {
			api.Configuration.SetGithubToken(nil)
		})
		api.Core.TopicController.WorkspaceController = api.Core.WorkspaceController
	}
}
