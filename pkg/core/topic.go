package core

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Topic struct {
	basePath string
	Name     string
}

func newTopic(basePath, name string) *Topic {
	return &Topic{
		basePath: basePath,
		Name:     name,
	}
}

func (t *Topic) Path() string {
	return filepath.Join(t.basePath, t.Name)
}

func (t *Topic) LastModified() time.Time {
	fi, err := os.Stat(t.Path())
	if err != nil {
		return time.Time{}
	}
	return fi.ModTime()
}

type Topics []*Topic

func (t Topics) Sorted() Topics {
	sort.Slice(t, func(i, j int) bool {
		return t[i].LastModified().After(t[j].LastModified())
	})
	return t
}
