package core

import (
	"errors"
	"mynav/pkg/system"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GianlucaP106/gotmux/gotmux"
)

// Api exposes all core api functions.
type Api struct {
	topics     *TopicRepository
	workspaces *WorkspaceRepository
	tmux       *gotmux.Tmux
	local      *LocalConfig
	global     *GlobalConfig
	updater    *updater
}

// Inits the Api.
func NewApi(dir string) (*Api, error) {
	// init global config
	global, err := newGlobalConfig()
	if err != nil {
		return nil, err
	}

	// init local config
	local, err := newLocalConfig(dir)
	if err != nil {
		return nil, err
	}

	// if no local config return nil
	if local == nil {
		return nil, nil
	}

	tmux, err := gotmux.DefaultTmux()
	if err != nil {
		return nil, err
	}

	topics := newTopicRepository(local.path)
	workspaces := newWorkspaceRepository(topics.All())

	api := &Api{}
	api.topics = topics
	api.workspaces = workspaces
	api.tmux = tmux
	api.local = local
	api.global = global
	api.updater = &updater{}
	return api, nil
}

// Updates mynav by running update script.
func (a *Api) UpdateMynav() error {
	return a.updater.UpdateMynav()
}

// Sets the 'update-asked' global configuration variable to now.
func (a *Api) AskUpdate() {
	a.global.SetUpdateAsked(time.Now())
}

// Returns if an update is available
func (a *Api) UpdateAvailable() (bool, string) {
	return a.updater.UpdateAvailable()
}

// Returns LocalConfig.
func (a *Api) LocalConfig() *LocalConfigData {
	return a.local.ConfigData()
}

// Returns GlobalConfig.
func (a *Api) GlobalConfig() *GlobalConfigData {
	return a.global.ConfigData()
}

// Creates a new topic.
func (a *Api) NewTopic(name string) (*Topic, error) {
	topic := newTopic(name, filepath.Join(a.local.path, name))
	if err := a.topics.Save(topic); err != nil {
		return nil, err
	}
	return topic, nil
}

// Finds topic by name.
func (a *Api) FindTopicByName(name string) *Topic {
	return a.topics.FindByName(name)
}

// Returns all topics.
func (a *Api) AllTopics() Topics {
	return a.topics.All()
}

// Returns topic count.
func (a *Api) TopicCount() int {
	return a.topics.Count()
}

// Deletes a topic.
func (a *Api) DeleteTopic(t *Topic) error {
	workspaces := a.workspaces.All().ByTopic(t)

	if err := a.topics.Delete(t); err != nil {
		return err
	}

	for _, w := range workspaces {
		a.KillSession(w)

		// remove straight from container because parent dir has been deleted.
		a.workspaces.container.Remove(w)
	}
	return nil
}

// Renames a topic.
func (a *Api) RenameTopic(t *Topic, name string) error {
	workspaces := a.workspaces.All().ByTopic(t)
	t.Name = name
	if err := a.topics.Save(t); err != nil {
		return err
	}

	for _, w := range workspaces {
		newPath := filepath.Join(w.Topic.path, w.Name)
		session := a.Session(w)
		if session != nil {
			session.Rename(newPath)
		}
		w.path = newPath
	}

	return nil
}

// Creates a new workspace.
func (a *Api) NewWorkspace(t *Topic, name string) (*Workspace, error) {
	name = strings.ReplaceAll(name, ".", "_")
	w := newWorkspace(t, name)
	if err := a.workspaces.Save(w); err != nil {
		return nil, err
	}

	a.SelectWorkspace(w)
	return w, nil
}

// Finds workspace by name.
func (a *Api) FindWorkspaceByShortPath(shortPath string) *Workspace {
	return a.workspaces.FindByShortPath(shortPath)
}

// Returns all workspaces.
func (a *Api) AllWorkspaces() Workspaces {
	return a.workspaces.All()
}

// Returns the workspace count.
func (a *Api) WorkspacesCount() int {
	return a.workspaces.Count()
}

// Deletes this workspace.
func (a *Api) DeleteWorkspace(w *Workspace) error {
	a.KillSession(w)
	selected := a.SelectedWorkspace()
	if selected == w {
		a.SelectWorkspace(nil)
	}

	return a.workspaces.Delete(w)
}

// Renames the workspace.
func (a *Api) RenameWorkspace(w *Workspace, name string) error {
	name = strings.ReplaceAll(name, ".", "_")

	// if there exists a workspace with this name in the same topic
	if len(a.AllWorkspaces().ByTopic(w.Topic).ByName(name)) > 0 {
		return errors.New("workspace with this name already exists")
	}

	// store session before rename
	s := a.Session(w)

	// rename and save
	w.Name = name
	if err := a.workspaces.Save(w); err != nil {
		return err
	}

	// rename session to new path
	if s != nil {
		s.Rename(w.Path())
	}

	return a.SelectWorkspace(w)
}

// Moves the workspace to a different topic.
func (a *Api) MoveWorkspace(w *Workspace, topic *Topic) error {
	if w.Topic.Name == topic.Name {
		return errors.New("workspace is already in this topic")
	}

	// if there exists a workspace with this name in the same topic
	if len(a.AllWorkspaces().ByTopic(topic).ByName(w.Name)) > 0 {
		return errors.New("workspace with this name already exists")
	}

	// get session if this workspace
	s := a.Session(w)

	// change topic and save
	w.Topic = topic
	err := a.workspaces.Save(w)
	if err != nil {
		return err
	}

	// rename session to new path
	if s != nil {
		s.Rename(w.path)
	}

	return a.SelectWorkspace(w)
}

// Gets the persisted selected workspace.
func (a *Api) SelectedWorkspace() *Workspace {
	workspaceShortPath := a.LocalConfig().SelectedWorkspace
	return a.FindWorkspaceByShortPath(workspaceShortPath)
}

// Sets the persisted selected workspace.
func (a *Api) SelectWorkspace(w *Workspace) error {
	set := a.local.SetSelectedWorkspace
	if w != nil {
		set(w.ShortPath())
	} else {
		set("")
	}
	return nil
}

// Opens neovim in the provided workspace.
func (a *Api) OpenWorkspaceNvim(w *Workspace) error {
	a.SelectWorkspace(w)
	return system.CommandWithRedirect("nvim", w.Path()).Run()
}

// Opens terminal in the provided workspace.
func (a *Api) OpenWorkspaceTerminal(w *Workspace) error {
	cmd, err := system.OpenTerminalCmd(w.Path())
	if err != nil {
		return err
	}
	return cmd.Run()
}

// Clones repo into workspace.
func (a *Api) CloneWorkspaceRepo(w *Workspace, url string) error {
	a.SelectWorkspace(w)
	return w.CloneRepo(url)
}

// Opens the workspace with either the command in settings or the default.
func (a *Api) OpenWorkspace(w *Workspace) error {
	if IsTmuxSession() {
		err := a.OpenWorkspaceNvim(w)
		if err != nil {
			return err
		}
	} else {
		err := a.OpenSession(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// Returns the preview for a given workspace (looks for its session).
func (a *Api) WorkspacePreview(w *Workspace) string {
	s := a.Session(w)
	if s == nil {
		return ""
	}

	ws, _ := s.ListWindows()
	if len(ws) == 0 {
		return ""
	}

	ps, _ := ws[0].ListPanes()
	if len(ps) == 0 {
		return ""
	}

	p, _ := ps[0].Capture()
	return p
}

// Wraps a workspace and tmux session together.
type Session struct {
	*gotmux.Session
	Workspace *Workspace
}

// Session map for quickly looking up session by its name.
type SessionMap map[string]*gotmux.Session

// Returns a session.
func newSession(s *gotmux.Session, w *Workspace) *Session {
	return &Session{
		Session:   s,
		Workspace: w,
	}
}

// Returns a Session by workspace using its path, nil if doesnt exist.
func (s SessionMap) Get(w *Workspace) *Session {
	session := s[w.Path()]
	if session == nil {
		return nil
	}

	return newSession(session, w)
}

// Returns the session associated to this workspace, nil if doesnt exist.
func (a *Api) Session(w *Workspace) *Session {
	s, _ := a.tmux.GetSessionByName(w.Path())
	if s == nil {
		return nil
	}
	return newSession(s, w)
}

// Creates and/or attaches to the workspace session.
func (a *Api) OpenSession(w *Workspace) error {
	// select the workspace
	a.SelectWorkspace(w)

	// attach to the exist existing if there is one
	existing := a.Session(w)
	if existing != nil {
		return existing.Attach()
	}

	// create a new session
	session, err := a.tmux.NewSession(&gotmux.SessionOptions{
		Name:           w.path,
		StartDirectory: w.Path(),
	})
	if err != nil {
		return err
	}

	return session.Attach()
}

// Kills the workspace session.
func (a *Api) KillSession(w *Workspace) error {
	session := a.Session(w)
	if session == nil {
		return errors.New("session does not exist")
	}
	return session.Kill()
}

// Returns the number of workspaces active workspace sessions.
func (a *Api) SessionCount() int {
	return len(a.AllSessions())
}

// Returns all workspace sessions.
func (a *Api) AllSessions() []*Session {
	// build map of sessions to lookup by name
	sMap := a.SessionMap()

	// get all workspaces and check which are associated to session
	out := make([]*Session, 0)
	for _, w := range a.AllWorkspaces() {
		s := sMap.Get(w)
		if s != nil {
			out = append(out, s)
		}
	}

	return out
}

// Returns a session map to allow look up by its name.
func (a *Api) SessionMap() SessionMap {
	sMap := make(SessionMap)
	sessions, _ := a.tmux.ListSessions()
	for _, s := range sessions {
		sMap[s.Name] = s
	}
	return sMap
}

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}
