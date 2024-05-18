package api

import "os"

type Api struct {
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

func (c *Api) GetSystemStats() (numTopics int, numWorkspaces int) {
	numTopics = c.GetTopicCount()
	numWorkspaces = c.GetWorkspaceCount()
	return
}

func (api *Api) InitConfiguration() {
	dir, _ := os.Getwd()
	api.InitConfig(dir)
	api.InitControllers()
}

func (api *Api) InitControllers() {
	if api.IsConfigInitialized {
		api.TopicController = NewTopicController(api.path)
		api.WorkspaceController = NewWorkspaceController(api.GetTopics(), api.GetWorkspaceStorePath())
		api.TopicController.WorkspaceController = api.WorkspaceController
	}
}
