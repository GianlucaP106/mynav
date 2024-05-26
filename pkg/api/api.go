package api

import "os"

type Api struct {
	*TmuxSessionController
	*WorkspaceController
	*TopicController
	*Configuration
}

func NewApi() *Api {
	api := &Api{}
	api.Configuration = NewConfiguration()
	api.InitControllers()
	return api
}

func (api *Api) GetSystemStats() (numTopics int, numWorkspaces int) {
	numTopics = api.GetTopicCount()
	numWorkspaces = api.GetWorkspaceCount()
	return
}

func (api *Api) InitConfiguration() {
	dir, _ := os.Getwd()
	api.InitConfig(dir)
	api.InitControllers()
}

func (api *Api) InitControllers() {
	if api.IsConfigInitialized {
		api.TmuxSessionController = NewTmuxSessionController()
		api.TopicController = NewTopicController(api.path, api.TmuxSessionController)
		api.WorkspaceController = NewWorkspaceController(api.GetTopics(), api.GetWorkspaceStorePath(), api.TmuxSessionController)
		api.TopicController.WorkspaceController = api.WorkspaceController
	}
}
