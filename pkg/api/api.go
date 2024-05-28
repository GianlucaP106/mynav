package api

import (
	"errors"
	"os"
)

type Api struct {
	*TmuxSessionController
	*WorkspaceController
	*TopicController
	*Configuration
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Configuration = NewConfiguration()

	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	if !api.IsConfigInitialized && cwd == home {
		return nil, errors.New("initializing mynav in the home directory is not supported")
	}
	api.InitControllers()
	return api, nil
}

func (api *Api) GetSystemStats() (numTopics int, numWorkspaces int) {
	numTopics = api.GetTopicCount()
	numWorkspaces = api.GetWorkspaceCount()
	return
}

func (api *Api) InitConfiguration() error {
	dir, _ := os.Getwd()
	if _, err := api.InitConfig(dir); err != nil {
		return errors.New("cannot initialize mynav in the home directory")
	}

	api.InitControllers()
	return nil
}

func (api *Api) InitControllers() {
	if api.IsConfigInitialized {
		api.TmuxSessionController = NewTmuxSessionController()
		api.TopicController = NewTopicController(api.path, api.TmuxSessionController)
		api.WorkspaceController = NewWorkspaceController(api.GetTopics(), api.GetWorkspaceStorePath(), api.TmuxSessionController)
		api.TopicController.WorkspaceController = api.WorkspaceController
	}
}
