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
