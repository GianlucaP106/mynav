package api

import (
	"mynav/pkg/utils"
	"sort"
	"time"
)

type Topic struct {
	Name string
	Path string
}

func newTopic(name string, path string) *Topic {
	return &Topic{
		Name: name,
		Path: path,
	}
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
	time, _ := utils.GetLastModifiedTime(t.Path)
	return time
}

func (t *Topic) GetLastModifiedTimeFormatted() string {
	time := t.GetLastModifiedTime().Format(TimeFormat())
	return time
}
