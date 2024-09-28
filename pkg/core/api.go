package core

import (
	"errors"
	"os"
)

type Api struct {
	Tmux                *TmuxController
	Topics              *TopicController
	Workspaces          *WorkspaceController
	GlobalConfiguration *GlobalConfiguration
	LocalConfiguration  *LocalConfiguration
	Github              *GithubController
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.LocalConfiguration = NewLocalConfiguration()
	gc, err := NewGlobalConfiguration()
	if err != nil {
		return nil, err
	}

	api.GlobalConfiguration = gc
	if api.LocalConfiguration.IsInitialized {
		if err := api.initControllers(); err != nil {
			return nil, err
		}
	} else {
		if err := api.initStandaloneControllers(); err != nil {
			return nil, err
		}
	}

	return api, nil
}

func (api *Api) InitConfiguration() error {
	dir, _ := os.Getwd()
	if _, err := api.LocalConfiguration.InitConfig(dir); err != nil {
		return errors.New("cannot initialize mynav in the home directory")
	}

	api.initControllers()
	return nil
}

func (api *Api) initStandaloneControllers() error {
	tmux, err := NewTmuxController()
	api.Tmux = tmux
	if err != nil {
		return err
	}

	api.Github = NewGithubController(api.GlobalConfiguration)
	return nil
}

func (api *Api) initControllers() error {
	if !api.LocalConfiguration.IsInitialized {
		return nil
	}

	err := api.initStandaloneControllers()
	if err != nil {
		return err
	}

	api.Topics = NewTopicController(api.LocalConfiguration, api.Tmux)
	api.Workspaces = NewWorkspaceController(
		api.Topics.GetTopics(),
		api.Tmux,
		api.GlobalConfiguration,
		api.LocalConfiguration,
	)
	api.Topics.workspaceController = api.Workspaces

	return nil
}
