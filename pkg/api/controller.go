package api

import "os"

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

func (c *Controller) CreateTopic(name string) error {
	if _, err := c.TopicManager.CreateTopic(name); err != nil {
		return err
	}

	return nil
}

func (c *Controller) IsConfigInitialized() bool {
	return c.Configuration.ConfigInitialized
}

func (c *Controller) InitConfiguration() {
	dir, _ := os.Getwd()
	c.Configuration.InitConfig(dir)

	c.initManagers()
}

func (c *Controller) GetTopic(idx int) *Topic {
	if idx >= len(c.TopicManager.Topics) || idx < 0 {
		return nil
	}
	return c.TopicManager.Topics[idx]
}

func (c *Controller) GetTopics() Topics {
	return c.TopicManager.Topics.Sorted()
}

func (c *Controller) DeleteTopic(topic *Topic) {
	if topic != nil {
		c.TopicManager.DeleteTopic(topic)
	}
}

func (c *Controller) GetTopicCount() int {
	return len(c.TopicManager.Topics)
}

func (c *Controller) GetWorkspaceCount() int {
	return len(c.WorkspaceManager.Workspaces)
}

func (c *Controller) GetWorkspacesByTopic(topic *Topic) Workspaces {
	out := make(Workspaces, 0)
	for _, workspace := range c.GetWorkspaces() {
		if workspace.Topic == topic {
			out = append(out, workspace)
		}
	}
	return out.Sorted()
}

func (c *Controller) GetWorkspacesByTopicCount(topic *Topic) int {
	return c.GetWorkspacesByTopic(topic).Len()
}

func (c *Controller) GetWorkspace(idx int) *Workspace {
	if idx >= len(c.WorkspaceManager.Workspaces) || idx < 0 {
		return nil
	}
	return c.WorkspaceManager.Workspaces[idx]
}

func (c *Controller) GetWorkspaces() Workspaces {
	return c.WorkspaceManager.Workspaces
}

func (c *Controller) CreateWorkspace(name string, repoUrl string, topic *Topic) error {
	if _, err := c.WorkspaceManager.CreateWorkspace(name, repoUrl, topic); err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteWorkspace(workspace *Workspace) {
	if workspace != nil {
		c.WorkspaceManager.DeleteWorkspace(workspace)
	}
}
