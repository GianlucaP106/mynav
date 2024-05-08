package controller

import (
	"mynav/pkg/core"
)

type Controller struct {
	fs *core.Filesystem
}

func NewController() *Controller {
	return &Controller{
		fs: core.NewFilesystem(),
	}
}

func (c *Controller) CreateTopic(name string) error {
	if _, err := c.fs.CreateTopic(name); err != nil {
		return err
	}

	return nil
}

func (c *Controller) IsConfigInitialized() bool {
	return c.fs.ConfigInitialized
}

func (c *Controller) CreateConfigurationFile() {
	c.fs.CreateConfigurationFile()
	c.fs.InitFilesystem()
}

func (c *Controller) GetTopic(idx int) *core.Topic {
	if idx >= len(c.fs.Topics) || idx < 0 {
		return nil
	}
	return c.fs.Topics[idx]
}

func (c *Controller) GetTopics() core.Topics {
	return c.fs.Topics.Sorted()
}

func (c *Controller) DeleteTopic(topic *core.Topic) {
	if topic != nil {
		c.fs.DeleteTopic(topic)
	}
}

func (c *Controller) GetTopicCount() int {
	return len(c.fs.Topics)
}

func (c *Controller) GetWorkspaceCount() int {
	return len(c.fs.Workspaces)
}

func (c *Controller) GetWorkspacesByTopic(topic *core.Topic) core.Workspaces {
	out := make(core.Workspaces, 0)
	for _, workspace := range c.GetWorkspaces() {
		if workspace.Topic == topic {
			out = append(out, workspace)
		}
	}
	return out.Sorted()
}

func (c *Controller) GetWorkspacesByTopicCount(topic *core.Topic) int {
	return c.GetWorkspacesByTopic(topic).Len()
}

func (c *Controller) GetWorkspace(idx int) *core.Workspace {
	if idx >= len(c.fs.Workspaces) || idx < 0 {
		return nil
	}
	return c.fs.Workspaces[idx]
}

func (c *Controller) GetWorkspaces() core.Workspaces {
	return c.fs.Workspaces
}

func (c *Controller) CreateWorkspace(name string, repoUrl string, topic *core.Topic) error {
	if _, err := c.fs.CreateWorkspace(name, repoUrl, topic); err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteWorkspace(workspace *core.Workspace) {
	if workspace != nil {
		c.fs.DeleteWorkspace(workspace)
	}
}
