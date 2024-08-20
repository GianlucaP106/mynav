package core

import (
	"os"
	"path/filepath"
)

type TopicController struct {
	TopicRepository     *TopicRepository
	WorkspaceController *WorkspaceController
	TmuxController      *TmuxController
	rootPath            string
}

func NewTopicController(rootPath string, tsc *TmuxController) *TopicController {
	tc := &TopicController{
		rootPath:       rootPath,
		TmuxController: tsc,
	}
	tc.TopicRepository = NewTopicRepository(rootPath)
	return tc
}

func (tc *TopicController) CreateTopic(name string) error {
	topic := newTopic(name, filepath.Join(tc.rootPath, name))
	tc.TopicRepository.Save(topic)
	return nil
}

func (tc *TopicController) GetTopics() Topics {
	return tc.TopicRepository.TopicContainer.All()
}

func (tc *TopicController) GetTopicCount() int {
	return tc.TopicRepository.TopicContainer.Size()
}

func (tc *TopicController) DeleteTopic(t *Topic) error {
	if err := tc.TopicRepository.Delete(t); err != nil {
		return err
	}

	tc.WorkspaceController.DeleteWorkspacesByTopic(t)
	return nil
}

func (tc *TopicController) RenameTopic(t *Topic, newName string) error {
	wr := tc.WorkspaceController.WorkspaceRepository

	newTopicPath := filepath.Join(filepath.Dir(t.Path), newName)
	if err := os.Rename(t.Path, newTopicPath); err != nil {
		return err
	}

	topicWorkspaces := tc.WorkspaceController.GetWorkspaces().FilterByTopic(t)

	for _, w := range topicWorkspaces {
		newWorkspacePath := filepath.Join(newTopicPath, w.Name)
		newShortPath := filepath.Join(newName, w.Name)

		wr.Container.Delete(w.ShortPath())
		wr.DeleteMetadata(w)

		if wr.Datasource.GetData().SelectedWorkspace == w.ShortPath() {
			wr.Datasource.GetData().SelectedWorkspace = newShortPath
		}

		if s := tc.TmuxController.GetTmuxSessionByName(w.Path); s != nil {
			tc.TmuxController.RenameTmuxSession(s, newWorkspacePath)
		}

		wr.Container.Set(w.ShortPath(), w)
		w.Path = newWorkspacePath

		data := wr.Datasource.GetData()
		data.Workspaces[newShortPath] = w.Metadata
		wr.Datasource.SaveData(data)

	}

	t.Name = newName
	t.Path = newTopicPath
	return nil
}
