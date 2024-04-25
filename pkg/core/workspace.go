package core

import (
	"log"
	"mynav/pkg/utils"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

type Workspace struct {
	Topic      *Topic
	Filesystem *Filesystem
	Name       string
	Path       string
	GitRemote  string
}

func newWorkspace(name string, topic *Topic, fs *Filesystem) *Workspace {
	wsPath := filepath.Join(fs.path, filepath.Join(topic.Name, name))
	ws := &Workspace{
		Name:  name,
		Topic: topic,
		Path:  wsPath,
	}
	ws.detectGitRemote()
	return ws
}

func (ws *Workspace) detectGitRemote() {
	gitPath := ws.Path + "/.git"
	if _, err := filepath.Abs(gitPath); err != nil {
		return
	}

	gitRemote, _ := utils.GitRemote(gitPath)
	ws.GitRemote = gitRemote
}

func (ws *Workspace) OpenWorkspace() {
	err := exec.Command("open", "-a", "warp", ws.Path).Run()
	if err != nil {
		log.Panicln(err)
	}
}

type Workspaces []*Workspace

func (w Workspaces) Len() int      { return len(w) }
func (w Workspaces) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w Workspaces) Less(i, j int) bool {
	return w[i].GetLastModifiedTime().After(w[j].GetLastModifiedTime())
}

func (w Workspaces) Sorted() Workspaces {
	sort.Sort(w)
	return w
}

func (w *Workspace) GetLastModifiedTime() time.Time {
	time, _ := utils.GetLastModifiedTime(w.Path)
	return time
}

func (w *Workspace) GetLastModifiedTimeFormatted() string {
	time := w.GetLastModifiedTime().Format(w.Filesystem.getTimeFormat())
	return time
}
