package api

import (
	"errors"
	"os"
)

type Api struct {
	*TmuxController
	*WorkspaceController
	*TopicController
	*PortController
	*Configuration
}

func NewApi() (*Api, error) {
	api := &Api{}
	api.Configuration = NewConfiguration()

	if api.IsConfigInitialized {
		api.InitControllers()
	} else {
		api.InitTmuxController()
	}

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

func (api *Api) InitTmuxController() {
	api.PortController = NewPortController()
	api.TmuxController = NewTmuxController(api.PortController)
	api.PortController.TmuxController = api.TmuxController
}

func (api *Api) InitControllers() {
	if api.IsConfigInitialized {
		api.InitTmuxController()
		api.TopicController = NewTopicController(api.path, api.TmuxController)
		api.WorkspaceController = NewWorkspaceController(
			api.GetTopics(),
			api.GetWorkspaceStorePath(),
			api.TmuxController,
			api.PortController,
		)

		api.TopicController.WorkspaceController = api.WorkspaceController
	}
}
