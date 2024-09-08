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
	IpcClient           *IpcClient
	IpcServer           *IpcServer
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Core = &Core{}
	api.LocalConfiguration = NewLocalConfiguration()
	gc, err := NewGlobalConfiguration()
	if err != nil {
		return nil, err
	}

	api.GlobalConfiguration = gc
	if api.LocalConfiguration.IsConfigInitialized {
		if err := api.InitControllers(); err != nil {
			return nil, err
		}
	} else {
		if err := api.InitStandaloneController(); err != nil {
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

	api.InitControllers()
	return nil
}

func (api *Api) InitStandaloneController() error {
	tmux, err := NewTmuxController()
	api.Tmux = tmux
	if err != nil {
		return err
	}

	api.Github = NewGithubController(api.GlobalConfiguration)
	return nil
}

func (api *Api) InitControllers() error {
	if api.LocalConfiguration.IsConfigInitialized {
		err := api.InitStandaloneController()
		if err != nil {
			return err
		}

		api.Core.TopicController = NewTopicController(api.LocalConfiguration.GetLocalConfigDir(), api.Tmux)
		api.Core.WorkspaceController = NewWorkspaceController(
			api.Core.GetTopics(),
			api.LocalConfiguration.GetWorkspaceStorePath(),
			api.Tmux,
		)

		api.Core.TopicController.WorkspaceController = api.Core.WorkspaceController
		api.Core.WorkspaceController.GlobalConfiguration = api.GlobalConfiguration
	}

	return nil
}

func (api *Api) InitIpc(runAction func(func())) {
	socketPath := api.LocalConfiguration.GetSocketPath()
	if IsParentAppInstance() {
		// funcs
		api.IpcServer = NewIpcServer(socketPath)
	}

	api.IpcClient = NewIpcClient(socketPath)
	api.Core.WorkspaceController.IpcClient = api.IpcClient
}
