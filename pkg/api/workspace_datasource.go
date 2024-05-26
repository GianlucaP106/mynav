package api

import "mynav/pkg/utils"

type WorkspaceDataSchema struct {
	Workspaces        map[string]*WorkspaceMetadata `json:"workspaces"`
	SelectedWorkspace string                        `json:"selected-workspace"`
}

type WorkspaceDatasource struct {
	Data      *WorkspaceDataSchema
	StorePath string
}

func NewWorkspaceDatasource(storePath string) *WorkspaceDatasource {
	ds := &WorkspaceDatasource{
		StorePath: storePath,
	}
	ds.LoadStore()
	return ds
}

func (wd *WorkspaceDatasource) SaveMetadata(w *Workspace) {
	wd.Data.Workspaces[w.ShortPath()] = w.Metadata
	wd.SaveStore()
}

func (wd *WorkspaceDatasource) DeleteMetadata(w *Workspace) {
	delete(wd.Data.Workspaces, w.ShortPath())
	wd.SaveStore()
}

func (wd *WorkspaceDatasource) GetMetadata(w *Workspace) *WorkspaceMetadata {
	return wd.Data.Workspaces[w.ShortPath()]
}

func (wd *WorkspaceDatasource) SetSelectedWorkspace(w *Workspace) {
	wd.Data.SelectedWorkspace = w.ShortPath()
}

func (wd *WorkspaceDatasource) Sync(w WorkspaceContainer) {
	for id, m := range wd.Data.Workspaces {
		if w.Get(id) == nil || m.Description == "" {
			delete(wd.Data.Workspaces, id)
		}
	}

	id := wd.Data.SelectedWorkspace
	if w.Get(id) == nil {
		wd.Data.SelectedWorkspace = ""
	}

	wd.SaveStore()
}

func (wd *WorkspaceDatasource) LoadStore() {
	wd.Data = utils.Load[WorkspaceDataSchema](wd.StorePath)
	if wd.Data == nil {
		wd.Data = &WorkspaceDataSchema{
			Workspaces:        map[string]*WorkspaceMetadata{},
			SelectedWorkspace: "",
		}
	}
}

func (wd *WorkspaceDatasource) SaveStore() {
	utils.Save(wd.Data, wd.StorePath)
}
