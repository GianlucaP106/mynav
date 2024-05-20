package api

import (
	"os"
	"path/filepath"
)

type TopicController struct {
	TopicRepoitory      *TopicRepository
	WorkspaceController *WorkspaceController
	rootPath            string
}

func NewTopicController(rootPath string) *TopicController {
	tc := &TopicController{
		rootPath: rootPath,
	}
	tc.TopicRepoitory = NewTopicRepository(rootPath)
	return tc
}

func (tc *TopicController) CreateTopic(name string) error {
	topic := newTopic(name, filepath.Join(tc.rootPath, name))
	tc.TopicRepoitory.Save(topic)
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

		wr.WorkspaceContainer.Delete(w)
		wr.WorkspaceDatasource.DeleteMetadata(w)

		if wr.WorkspaceDatasource.Data.SelectedWorkspace == w.ShortPath() {
			wr.WorkspaceDatasource.Data.SelectedWorkspace = newShortPath
		}

		wr.WorkspaceContainer[newShortPath] = w
		w.Path = newWorkspacePath

		if w.Metadata.TmuxSession != nil {
			tc.WorkspaceController.TmuxCommunicator.RenameSession(w.Metadata.TmuxSession, newWorkspacePath)
			w.Metadata.TmuxSession.Name = newWorkspacePath
		}

		wr.WorkspaceDatasource.Data.Workspaces[newShortPath] = w.Metadata
		wr.WorkspaceDatasource.SaveStore()

	}

	t.Name = newName
	return nil
}
