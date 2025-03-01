package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/GianlucaP106/mynav/pkg/system"
)

// Topic
type Topic struct {
	Name string
	path string

	// TODO: parent topic
}

func newTopic(name string, path string) *Topic {
	return &Topic{
		Name: name,
		path: path,
	}
}

func (t *Topic) LastModifiedTime() time.Time {
	time, _ := system.GetLastModifiedTime(t.path)
	return time
}

func (t *Topic) LastModifiedTimeFormatted() string {
	time := t.LastModifiedTime().Format(system.TimeFormat())
	return time
}

// Topics is a collection of Topic.
type Topics []*Topic

func (t Topics) Len() int { return len(t) }

func (t Topics) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t Topics) Less(i, j int) bool {
	return t[i].LastModifiedTime().After(t[j].LastModifiedTime())
}

func (t Topics) Sorted() Topics {
	sort.Sort(t)
	return t
}

func (t Topics) ByNameContaining(s string) Topics {
	if s == "" {
		return t
	}

	filtered := Topics{}
	for _, topic := range t {
		if strings.Contains(topic.Name, s) {
			filtered = append(filtered, topic)
		}
	}
	return filtered
}

func (t Topics) GetTopic(idx int) *Topic {
	if idx >= len(t) || idx < 0 {
		return nil
	}
	return t[idx]
}

// TopicRepository exposes basic crud on topics.
type TopicRepository struct {
	container *Container[Topic]
}

func newTopicRepository(path string) *TopicRepository {
	tr := &TopicRepository{}
	tr.load(path)
	return tr
}

func (tr *TopicRepository) load(rootPath string) {
	tc := newContainer[Topic]()
	tr.container = tc
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
	oldName := filepath.Base(t.path)
	// if this topic doesnt exist, we create a dir
	if !tr.container.Contains(oldName) {
		if err := system.CreateDir(t.path); err != nil {
			return err
		}
	}

	// if the name changed (based on its path) we rename the dir
	if t.Name != filepath.Base(t.path) {
		newPath := filepath.Join(filepath.Dir(t.path), t.Name)
		if err := os.Rename(t.path, newPath); err != nil {
			return err
		}
		tr.container.Remove(oldName)
		t.path = newPath
	}

	// save it to the container
	tr.container.Set(t.Name, t)
	return nil
}

func (tr *TopicRepository) Delete(t *Topic) error {
	if err := os.RemoveAll(t.path); err != nil {
		return err
	}

	tr.container.Remove(t.Name)
	return nil
}

func (tr *TopicRepository) FindByName(name string) *Topic {
	return tr.container.Get(name)
}

func (tr *TopicRepository) All() Topics {
	return tr.container.All()
}

func (tr *TopicRepository) Count() int {
	return tr.container.Size()
}
