package api

import (
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

type TopicRepoitory struct {
	TopicContainer TopicContainer
}

func NewTopicRepository(rootPath string) *TopicRepoitory {
	tr := &TopicRepoitory{}
	tr.LoadContainer(rootPath)
	return tr
}

func (tr *TopicRepoitory) LoadContainer(rootPath string) {
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

func (tr *TopicRepoitory) Save(t *Topic) error {
	existing := tr.TopicContainer.Get(t.Name)
	if existing == nil {
		if err := utils.CreateDir(t.Path); err != nil {
			return err
		}
	}

	tr.TopicContainer.Set(t)
	return nil
}

func (tr *TopicRepoitory) Delete(t *Topic) error {
	if err := os.RemoveAll(t.Path); err != nil {
		return err
	}

	tr.TopicContainer.Delete(t)
	return nil
}

func (tr *TopicRepoitory) Find(name string) *Topic {
	return tr.TopicContainer.Get(name)
}
