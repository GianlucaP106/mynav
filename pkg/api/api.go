package api

import (
	"errors"
	"mynav/pkg/configuration"
	"mynav/pkg/core"
	"mynav/pkg/github"
	"mynav/pkg/system"
	"os"
)

type Api struct {
	Tmux          *core.TmuxController
	Core          *core.Core
	Configuration *configuration.Configuration
	Github        *github.GithubController
	Port          *system.PortController
	Proc          *system.ProcessController
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Core = &core.Core{}
	api.Configuration = &configuration.Configuration{}
	api.Configuration.LocalConfiguration = configuration.NewLocalConfiguration()
	api.Configuration.GlobalConfiguration = configuration.NewGlobalConfiguration()

	if api.Configuration.IsConfigInitialized {
		api.InitControllers()
	} else {
		api.InitStandaloneController()
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

func (api *Api) InitStandaloneController() {
	api.Proc = system.NewProcessController()
	api.Port = system.NewPortController(api.Proc)
	api.Tmux = core.NewTmuxController(api.Port, api.Proc)

	// TODO: move config to seperate module
	api.Github = github.NewGithubController(api.Configuration.GetGithubToken(), func(gat *github.GithubAuthenticationToken) {
		api.Configuration.SetGithubToken(gat)
	}, func() {
		api.Configuration.SetGithubToken(nil)
	})
}

func (api *Api) InitControllers() {
	if api.Configuration.IsConfigInitialized {
		api.InitStandaloneController()
		api.Core.TopicController = core.NewTopicController(api.Configuration.GetLocalConfigDir(), api.Tmux)
		api.Core.WorkspaceController = core.NewWorkspaceController(
			api.Core.GetTopics(),
			api.Configuration.GetWorkspaceStorePath(),
			api.Tmux,
		)
		api.Core.TopicController.WorkspaceController = api.Core.WorkspaceController
	}
}
