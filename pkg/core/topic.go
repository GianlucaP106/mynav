package core

import (
	"mynav/pkg/utils"
	"path/filepath"
	"sort"
	"time"
)

type Topic struct {
	Filesystem *Filesystem
	Name       string
}

func newTopic(name string, fs *Filesystem) *Topic {
	return &Topic{
		Name:       name,
		Filesystem: fs,
	}
}

func (t *Topic) GetPath() string {
	return filepath.Join(t.Filesystem.path, t.Name)
}

type Topics []*Topic

func (t Topics) Len() int      { return len(t) }
func (t Topics) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Topics) Less(i, j int) bool {
	return t[i].GetLastModifiedTime().After(t[j].GetLastModifiedTime())
}

func (t Topics) Sorted() Topics {
	sort.Sort(t)
	return t
}

func (t *Topic) GetLastModifiedTime() time.Time {
	time, _ := utils.GetLastModifiedTime(t.GetPath())
	return time
}

func (t *Topic) GetLastModifiedTimeFormatted() string {
	time := t.GetLastModifiedTime().Format(t.Filesystem.getTimeFormat())
	return time
}
