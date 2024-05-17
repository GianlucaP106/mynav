package api

import "path/filepath"

type TopicController struct {
	TopicRepoitory      *TopicRepoitory
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
