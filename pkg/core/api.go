package core

import (
	"os"
	"path/filepath"

	"github.com/GianlucaP106/gotmux/gotmux"
)

// API exposes all core api functions.
type API struct {
	fs      *Filesystem
	tmux    *gotmux.Tmux
	local   *LocalConfig
	global  *GlobalConfig
	updater *updater
}

// Inits the Api.
func NewApi(dir string) (*API, error) {
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

	c := newFilesystem(local.path)

	api := &API{}
	api.fs = c
	api.tmux = tmux
	api.local = local
	api.global = global
	api.updater = &updater{}
	return api, nil
}

// Returns if an update is available
func (a *API) UpdateAvailable() (bool, string) {
	return a.updater.UpdateAvailable()
}

// Creates a new topic.
func (a *API) NewTopic(name string) (*Topic, error) {
	return a.fs.CreateTopic(name)
}

// Returns all topics.
func (a *API) Topics() Topics {
	return a.fs.Topics()
}

// Returns topic count.
func (a *API) TopicCount() int {
	return a.fs.TopicsCount()
}

// Deletes a topic.
func (a *API) DeleteTopic(t *Topic) error {
	for _, w := range a.Workspaces(t) {
		if s := a.Session(w); s != nil {
			s.Kill()
		}
	}
	a.SelectWorkspace(nil)
	return a.fs.DeleteTopic(t)
}

// Renames a topic.
func (a *API) RenameTopic(t *Topic, name string) error {
	// store topic path for session rename
	oldTopicPath := t.Path()

	if err := a.fs.RenameTopic(t, name); err != nil {
		return err
	}

	// rename all sessions
	for _, w := range a.Workspaces(t) {
		session, err := a.tmux.GetSessionByName(filepath.Join(oldTopicPath, w.Name))
		if err != nil {
			return err
		}

		if session != nil {
			if err := session.Rename(w.Path()); err != nil {
				return err
			}
		}
	}

	return nil
}

// Creates a new workspace.
func (a *API) NewWorkspace(t *Topic, name string) (*Workspace, error) {
	w, err := a.fs.CreateWorkspace(t, name)
	if err != nil {
		return nil, err
	}
	a.SelectWorkspace(w)
	return w, nil
}

// Returns the workspaces for this topic.
func (a *API) Workspaces(t *Topic) Workspaces {
	return a.fs.Workspaces(t)
}

// Returns all workspaces.
func (a *API) AllWorkspaces() Workspaces {
	return a.fs.AllWorkspaces()
}

// Returns the workspace count.
func (a *API) WorkspacesCount() int {
	return a.fs.WorkspacesCount()
}

// Deletes this workspace.
func (a *API) DeleteWorkspace(w *Workspace) error {
	if s := a.Session(w); s != nil {
		s.Kill()
	}
	selected := a.SelectedWorkspace()
	if selected == w {
		a.SelectWorkspace(nil)
	}

	return a.fs.DeleteWorkspace(w)
}

// Renames the workspace.
func (a *API) RenameWorkspace(w *Workspace, name string) error {
	s := a.Session(w)

	if err := a.fs.RenameWorkspace(w, name); err != nil {
		return err
	}

	// rename session to new path
	if s != nil {
		s.Rename(w.Path())
	}

	return a.SelectWorkspace(w)
}

// Moves the workspace to a different topic.
func (a *API) MoveWorkspace(w *Workspace, topic *Topic) error {
	s := a.Session(w)
	if err := a.fs.MoveWorkspace(w, topic); err != nil {
		return err
	}

	// rename session to new path
	if s != nil {
		s.Rename(w.Path())
	}

	return a.SelectWorkspace(w)
}

// Gets the persisted selected workspace.
func (a *API) SelectedWorkspace() *Workspace {
	lcd := a.local.ConfigData()
	return a.Workspace(lcd.SelectedWorkspace)
}

// Sets the persisted selected workspace.
func (a *API) SelectWorkspace(w *Workspace) error {
	if w != nil {
		a.local.SetSelectedWorkspace(w.ShortPath())
	} else {
		a.local.SetSelectedWorkspace("")
	}
	return nil
}

// Returns a Workspace object if short path is valid.
func (s *API) Workspace(shortPath string) *Workspace {
	return s.fs.Workspace(shortPath)
}

// Clones repo into workspace.
func (a *API) CloneWorkspaceRepo(w *Workspace, url string) error {
	a.SelectWorkspace(w)
	return w.CloneRepo(url)
}

// Wraps a workspace and tmux session together.
type Session struct {
	*gotmux.Session
	Workspace *Workspace
}

func (s *Session) DisplayName() string {
	if s.Workspace == nil {
		return s.Name
	}

	return s.Workspace.ShortPath()
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
	session := s[w.TmuxName()]
	if session == nil {
		return nil
	}

	return newSession(session, w)
}

// Returns the session associated to this workspace, nil if doesnt exist.
func (a *API) Session(w *Workspace) *Session {
	if w == nil {
		return nil
	}

	s, _ := a.tmux.GetSessionByName(w.TmuxName())
	if s == nil {
		return nil
	}
	return newSession(s, w)
}

// Creates and/or attaches to the workspace session.
func (a *API) OpenWorkspace(w *Workspace) error {
	// select the workspace
	a.SelectWorkspace(w)

	// attach to the exist existing if there is one
	existing := a.Session(w)
	if existing != nil {
		return existing.Attach()
	}

	// create a new session
	path := w.Path()
	name := w.TmuxName()
	session, err := a.tmux.NewSession(&gotmux.SessionOptions{
		Name:           name,
		StartDirectory: path,
	})
	if err != nil {
		return err
	}

	return session.Attach()
}

func (a *API) NewSession(name string) (*Session, error) {
	s, err := a.tmux.NewSession(&gotmux.SessionOptions{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return newSession(s, nil), nil
}

func (a *API) SessionByName(name string) (*Session, error) {
	session, err := a.tmux.Session(name)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, nil
	}

	return newSession(session, nil), nil
}

// Returns the number of workspaces active workspace sessions.
func (a *API) SessionCount() int {
	return len(a.AllSessions())
}

// Returns all workspace sessions.
func (a *API) AllSessions() []*Session {
	// build map of workspaces accessible by tmux name
	wMap := map[string]*Workspace{}
	for _, w := range a.AllWorkspaces() {
		wMap[w.TmuxName()] = w
	}

	sessions, _ := a.tmux.ListSessions()

	out := make([]*Session, len(sessions))
	for i, s := range sessions {
		associatedWorkspace := wMap[s.Name]
		out[i] = newSession(s, associatedWorkspace)
	}

	return out
}

// Returns a session map to allow look up by its name.
func (a *API) SessionMap() SessionMap {
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
