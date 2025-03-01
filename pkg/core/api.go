package core

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/GianlucaP106/gotmux/gotmux"
)

// API exposes all core api functions.
type API struct {
	container *Container
	tmux      *gotmux.Tmux
	local     *LocalConfig
	global    *GlobalConfig
	updater   *updater
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

	c := newContainer(local.path)

	api := &API{}
	api.container = c
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
	return a.container.CreateTopic(name)
}

// Returns all topics.
func (a *API) AllTopics() Topics {
	return a.container.AllTopics()
}

// Returns topic count.
func (a *API) TopicCount() int {
	return a.container.TopicsCount()
}

// Deletes a topic.
func (a *API) DeleteTopic(t *Topic) error {
	for _, w := range t.workspaces {
		a.KillSession(w)
	}
	// TODO: selected
	return a.container.DeleteTopic(t)
}

// Renames a topic.
func (a *API) RenameTopic(t *Topic, name string) error {
	// store topic path for session rename
	oldTopicPath := t.Path()

	if err := a.container.RenameTopic(t, name); err != nil {
		return err
	}

	// rename all sessions
	for _, w := range t.workspaces {
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
	w, err := a.container.CreateWorkspace(t, name)
	if err != nil {
		return nil, err
	}
	a.SelectWorkspace(w)
	return w, nil
}

// Returns all workspaces.
func (a *API) AllWorkspaces() Workspaces {
	return a.container.AllWorkspaces()
}

// Returns the workspace count.
func (a *API) WorkspacesCount() int {
	return a.container.WorkspacesCount()
}

// Deletes this workspace.
func (a *API) DeleteWorkspace(w *Workspace) error {
	a.KillSession(w)
	selected := a.SelectedWorkspace()
	if selected == w {
		a.SelectWorkspace(nil)
	}

	return a.container.DeleteWorkspace(w)
}

// Renames the workspace.
func (a *API) RenameWorkspace(w *Workspace, name string) error {
	s := a.Session(w)

	if err := a.container.RenameWorkspace(w, name); err != nil {
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
	if err := a.container.MoveWorkspace(w, topic); err != nil {
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
	return a.FindWorkspace(lcd.SelectedWorkspace)
}

// Sets the persisted selected workspace.
func (a *API) SelectWorkspace(w *Workspace) error {
	set := a.local.SetSelectedWorkspace
	if w != nil {
		set(w.ShortPath())
	} else {
		set("")
	}
	return nil
}

func (s *API) FindWorkspace(shortPath string) *Workspace {
	topicName, workspaceName := filepath.Dir(shortPath), filepath.Base(shortPath)
	topic := s.container.topics[topicName]

	if topic == nil {
		return nil
	}

	return topic.workspaces[workspaceName]
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
func (a *API) Session(w *Workspace) *Session {
	s, _ := a.tmux.GetSessionByName(w.Path())
	if s == nil {
		return nil
	}
	return newSession(s, w)
}

// Creates and/or attaches to the workspace session.
func (a *API) OpenSession(w *Workspace) error {
	// select the workspace
	a.SelectWorkspace(w)

	// attach to the exist existing if there is one
	existing := a.Session(w)
	if existing != nil {
		return existing.Attach()
	}

	// create a new session
	p := w.Path()
	session, err := a.tmux.NewSession(&gotmux.SessionOptions{
		Name:           p,
		StartDirectory: p,
	})
	if err != nil {
		return err
	}

	return session.Attach()
}

// Kills the workspace session.
func (a *API) KillSession(w *Workspace) error {
	session := a.Session(w)
	if session == nil {
		return errors.New("session does not exist")
	}
	return session.Kill()
}

// Returns the number of workspaces active workspace sessions.
func (a *API) SessionCount() int {
	return len(a.AllSessions())
}

// Returns all workspace sessions.
func (a *API) AllSessions() []*Session {
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
