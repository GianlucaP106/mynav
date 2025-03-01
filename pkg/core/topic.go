package core

import (
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/GianlucaP106/mynav/pkg/system"
)

// Topic
type Topic struct {
	basePath   string
	Name       string
	workspaces map[string]*Workspace
	mu         sync.RWMutex
}

func newTopic(basePath, name string) *Topic {
	return &Topic{
		basePath:   basePath,
		Name:       name,
		workspaces: make(map[string]*Workspace),
	}
}

func (t *Topic) Path() string {
	return filepath.Join(t.basePath, t.Name)
}

func (t *Topic) LastModifiedTime() time.Time {
	time, _ := system.GetLastModifiedTime(t.Path())
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
