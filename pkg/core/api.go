package core

import (
	"errors"
	"os"
)

type Core struct {
	*TopicController
	*WorkspaceController
}

type Api struct {
	Tmux                *TmuxController
	Core                *Core
	GlobalConfiguration *GlobalConfiguration
	LocalConfiguration  *LocalConfiguration
	Github              *GithubController
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Core = &Core{}
	api.LocalConfiguration = NewLocalConfiguration()
	api.GlobalConfiguration = NewGlobalConfiguration()

	if api.LocalConfiguration.IsConfigInitialized {
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
	if _, err := api.LocalConfiguration.InitConfig(dir); err != nil {
		return errors.New("cannot initialize mynav in the home directory")
	}

	api.InitControllers()
	return nil
}

func (api *Api) InitStandaloneController() {
	api.Tmux = NewTmuxController()
	api.Github = NewGithubController(api.GlobalConfiguration)
}

func (api *Api) InitControllers() {
	if api.LocalConfiguration.IsConfigInitialized {
		api.InitStandaloneController()
		api.Core.TopicController = NewTopicController(api.LocalConfiguration.GetLocalConfigDir(), api.Tmux)
		api.Core.WorkspaceController = NewWorkspaceController(
			api.Core.GetTopics(),
			api.LocalConfiguration.GetWorkspaceStorePath(),
			api.Tmux,
		)
		api.Core.TopicController.WorkspaceController = api.Core.WorkspaceController
	}
}
