package core

import (
	"mynav/pkg/persistence"
	"mynav/pkg/system"
	"os"
	"path/filepath"
)

type TopicRepository struct {
	TopicContainer *persistence.Container[Topic]
}

func NewTopicRepository(rootPath string) *TopicRepository {
	tr := &TopicRepository{}
	tr.LoadContainer(rootPath)
	return tr
}

func (tr *TopicRepository) LoadContainer(rootPath string) {
	tc := persistence.NewContainer[Topic]()
	tr.TopicContainer = tc
	for _, topicDirEntry := range system.GetDirEntries(rootPath) {
		if !topicDirEntry.IsDir() || topicDirEntry.Name() == ".mynav" {
			continue
		}

		topicName := topicDirEntry.Name()
		topic := newTopic(topicName, filepath.Join(rootPath, topicName))
		tc.Set(topic.Name, topic)
	}
}

func (tr *TopicRepository) Save(t *Topic) error {
	existing := tr.TopicContainer.Get(t.Name)
	if existing == nil {
		if err := system.CreateDir(t.Path); err != nil {
			return err
		}
	}

	tr.TopicContainer.Set(t.Name, t)
	return nil
}

func (tr *TopicRepository) Delete(t *Topic) error {
	if err := os.RemoveAll(t.Path); err != nil {
		return err
	}

	tr.TopicContainer.Delete(t.Name)
	return nil
}

func (tr *TopicRepository) Find(name string) *Topic {
	return tr.TopicContainer.Get(name)
}
