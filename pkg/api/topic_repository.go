package api

import (
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

type TopicRepository struct {
	TopicContainer TopicContainer
}

func NewTopicRepository(rootPath string) *TopicRepository {
	tr := &TopicRepository{}
	tr.LoadContainer(rootPath)
	return tr
}

func (tr *TopicRepository) LoadContainer(rootPath string) {
	tc := NewTopicContainer()
	tr.TopicContainer = tc
	for _, topicDirEntry := range utils.GetDirEntries(rootPath) {
		if !topicDirEntry.IsDir() || topicDirEntry.Name() == ".mynav" {
			continue
		}

		topicName := topicDirEntry.Name()
		topic := newTopic(topicName, filepath.Join(rootPath, topicName))
		tc.Set(topic)
	}
}

func (tr *TopicRepository) Save(t *Topic) error {
	existing := tr.TopicContainer.Get(t.Name)
	if existing == nil {
		if err := utils.CreateDir(t.Path); err != nil {
			return err
		}
	}

	tr.TopicContainer.Set(t)
	return nil
}

func (tr *TopicRepository) Delete(t *Topic) error {
	if err := os.RemoveAll(t.Path); err != nil {
		return err
	}

	tr.TopicContainer.Delete(t)
	return nil
}

func (tr *TopicRepository) Find(name string) *Topic {
	return tr.TopicContainer.Get(name)
}
