package core

type WorkspaceContainer map[string]*Workspace

func NewWorkspaceContainer() WorkspaceContainer {
	return make(WorkspaceContainer)
}

func (wc WorkspaceContainer) Get(id string) *Workspace {
	return wc[id]
}

func (wc WorkspaceContainer) Set(w *Workspace) {
	wc[w.ShortPath()] = w
}

func (wc WorkspaceContainer) Delete(w *Workspace) {
	delete(wc, w.ShortPath())
}

func (wc WorkspaceContainer) ToList() Workspaces {
	out := make(Workspaces, 0)
	for _, w := range wc {
		out = append(out, w)
	}
	return out
}
