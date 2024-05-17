package api

type Api struct {
	*WorkspaceController
	*TopicController
	*Configuration
}

func NewApi() *Api {
	api := &Api{}
	api.Configuration = NewConfiguration()
	api.TopicController = NewTopicController(api.path)
	api.WorkspaceController = NewWorkspaceController(api.GetTopics(), api.GetWorkspaceStorePath())
	api.TopicController.WorkspaceController = api.WorkspaceController
	return api
}

func (c *Api) GetSystemStats() (numTopics int, numWorkspaces int) {
	numTopics = c.GetTopicCount()
	numWorkspaces = c.GetWorkspaceCount()
	return
}
