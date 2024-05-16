package api

import (
	"os"
)

type Controller struct {
	WorkspaceManager *WorkspaceManager
	TopicManager     *TopicManager
	Configuration    *Configuration
}

func NewController() *Controller {
	c := &Controller{}
	c.Configuration = &Configuration{}
	c.Configuration.DetectConfig()
	// utils.InitLogger(filepath.Join(c.Configuration.path, "debug.log"))
	c.initManagers()
	return c
}

func (c *Controller) initManagers() {
	if c.Configuration.ConfigInitialized {
		c.TopicManager = newTopicManager(c)
		c.WorkspaceManager = newWorkspaceManager(c)
	}
}

func (c *Controller) InitConfiguration() {
	dir, _ := os.Getwd()
	c.Configuration.InitConfig(dir)
	c.initManagers()
}

func (c *Controller) GetSystemStats() (numTopics int, numWorkspaces int) {
	numTopics = c.TopicManager.Topics.Len()
	numWorkspaces = c.WorkspaceManager.Workspaces.Len()
	return
}
