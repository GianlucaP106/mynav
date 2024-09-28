package core

import (
	"os"
	"path/filepath"
)

type TopicController struct {
	topicRepository     *TopicRepository
	workspaceController *WorkspaceController
	tmuxController      *TmuxController
	localConfiguration  *LocalConfiguration
}

func NewTopicController(lc *LocalConfiguration, tsc *TmuxController) *TopicController {
	tc := &TopicController{
		localConfiguration: lc,
		tmuxController:     tsc,
	}
	tc.topicRepository = NewTopicRepository(lc.path)
	return tc
}

func (tc *TopicController) CreateTopic(name string) error {
	topic := newTopic(name, filepath.Join(tc.localConfiguration.path, name))
	tc.topicRepository.Save(topic)
	return nil
}

func (tc *TopicController) GetTopics() Topics {
	return tc.topicRepository.container.All()
}

func (tc *TopicController) GetTopicCount() int {
	return tc.topicRepository.container.Size()
}

func (tc *TopicController) DeleteTopic(t *Topic) error {
	if err := tc.topicRepository.Delete(t); err != nil {
		return err
	}

	tc.workspaceController.DeleteWorkspacesByTopic(t)
	return nil
}

func (tc *TopicController) RenameTopic(t *Topic, newName string) error {
	wr := tc.workspaceController.workspaceRepository

	newTopicPath := filepath.Join(filepath.Dir(t.Path), newName)
	if err := os.Rename(t.Path, newTopicPath); err != nil {
		return err
	}

	topicWorkspaces := tc.workspaceController.GetWorkspaces().FilterByTopic(t)

	for _, w := range topicWorkspaces {
		newWorkspacePath := filepath.Join(newTopicPath, w.Name)
		newShortPath := filepath.Join(newName, w.Name)

		wr.container.Delete(w.ShortPath())
		wr.DeleteMetadata(w)

		if wr.datasource.GetData().SelectedWorkspace == w.ShortPath() {
			wr.datasource.GetData().SelectedWorkspace = newShortPath
		}

		if s := tc.tmuxController.GetTmuxSessionByName(w.Path); s != nil {
			tc.tmuxController.RenameTmuxSession(s, newWorkspacePath)
		}

		wr.container.Set(w.ShortPath(), w)
		w.Path = newWorkspacePath

		data := wr.datasource.GetData()
		data.Workspaces[newShortPath] = w.Metadata
		wr.datasource.SaveData(data)

	}

	t.Name = newName
	t.Path = newTopicPath
	return nil
}
