package core

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/tmux"
	"os"
	"path/filepath"
)

type TopicController struct {
	TopicRepoitory      *TopicRepository
	WorkspaceController *WorkspaceController
	TmuxController      *tmux.TmuxController
	rootPath            string
}

func NewTopicController(rootPath string, tsc *tmux.TmuxController) *TopicController {
	tc := &TopicController{
		rootPath:       rootPath,
		TmuxController: tsc,
	}
	tc.TopicRepoitory = NewTopicRepository(rootPath)
	return tc
}

func (tc *TopicController) CreateTopic(name string) error {
	topic := newTopic(name, filepath.Join(tc.rootPath, name))
	tc.TopicRepoitory.Save(topic)
	events.EmitEvent(constants.TopicChangeEventName)
	return nil
}

func (tc *TopicController) GetTopics() Topics {
	return tc.TopicRepoitory.TopicContainer.ToList()
}

func (tc *TopicController) GetTopicCount() int {
	return len(tc.TopicRepoitory.TopicContainer)
}

func (tc *TopicController) DeleteTopic(t *Topic) error {
	if err := tc.TopicRepoitory.Delete(t); err != nil {
		return err
	}

	tc.WorkspaceController.DeleteWorkspacesByTopic(t)
	events.EmitEvent(constants.TopicChangeEventName)
	return nil
}

func (tc *TopicController) RenameTopic(t *Topic, newName string) error {
	wr := tc.WorkspaceController.WorkspaceRepository

	newTopicPath := filepath.Join(filepath.Dir(t.Path), newName)
	if err := os.Rename(t.Path, newTopicPath); err != nil {
		return err
	}

	topicWorkspaces := tc.WorkspaceController.GetWorkspaces().ByTopic(t)

	for _, w := range topicWorkspaces {
		newWorkspacePath := filepath.Join(newTopicPath, w.Name)
		newShortPath := filepath.Join(newName, w.Name)

		wr.Container.Delete(w)
		wr.DeleteMetadata(w)

		if wr.Datasource.Data.SelectedWorkspace == w.ShortPath() {
			wr.Datasource.Data.SelectedWorkspace = newShortPath
		}

		if s := tc.TmuxController.GetTmuxSessionByName(w.Path); s != nil {
			tc.TmuxController.RenameTmuxSession(s, newWorkspacePath)
		}

		wr.Container[newShortPath] = w
		w.Path = newWorkspacePath

		wr.Datasource.Data.Workspaces[newShortPath] = w.Metadata
		wr.Datasource.SaveData()

	}

	t.Name = newName
	t.Path = newTopicPath
	events.EmitEvent(constants.TopicChangeEventName)
	events.EmitEvent(constants.WorkspaceChangeEventName)
	return nil
}
