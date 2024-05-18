package api

import "mynav/pkg/utils"

type WorkspaceDataSchema struct {
	Workspaces        map[string]*WorkspaceMetadata `json:"workspaces"`
	SelectedWorkspace string                        `json:"selected-workspace"`
}

type WorkspaceDatasource struct {
	WorkspaceStoreSchema *WorkspaceDataSchema
	StorePath            string
}

func NewWorkspaceDatasource(storePath string) *WorkspaceDatasource {
	ds := &WorkspaceDatasource{
		StorePath: storePath,
	}
	ds.LoadStore()
	return ds
}

func (wd *WorkspaceDatasource) SaveMetadata(w *Workspace) {
	wd.WorkspaceStoreSchema.Workspaces[w.ShortPath()] = w.Metadata
	wd.SaveStore()
}

func (wd *WorkspaceDatasource) DeleteMetadata(w *Workspace) {
	delete(wd.WorkspaceStoreSchema.Workspaces, w.ShortPath())
	wd.SaveStore()
}

func (wd *WorkspaceDatasource) GetMetadata(w *Workspace) *WorkspaceMetadata {
	return wd.WorkspaceStoreSchema.Workspaces[w.ShortPath()]
}

func (wd *WorkspaceDatasource) SetSelectedWorkspace(w *Workspace) {
	wd.WorkspaceStoreSchema.SelectedWorkspace = w.ShortPath()
	wd.SaveStore()
}

func (wd *WorkspaceDatasource) Sync(w WorkspaceContainer) {
	for id := range wd.WorkspaceStoreSchema.Workspaces {
		if w.Get(id) == nil {
			delete(wd.WorkspaceStoreSchema.Workspaces, id)
		}
	}

	id := wd.WorkspaceStoreSchema.SelectedWorkspace
	if w.Get(id) == nil {
		wd.WorkspaceStoreSchema.SelectedWorkspace = ""
	}

	wd.SaveStore()
}

func (wd *WorkspaceDatasource) LoadStore() {
	store := utils.Load[WorkspaceDataSchema](wd.StorePath)
	if store != nil {
		wd.WorkspaceStoreSchema = store
		return
	}

	wd.WorkspaceStoreSchema = &WorkspaceDataSchema{
		Workspaces:        map[string]*WorkspaceMetadata{},
		SelectedWorkspace: "",
	}

	wd.SaveStore()
}

func (wd *WorkspaceDatasource) SaveStore() {
	utils.Save(wd.WorkspaceStoreSchema, wd.StorePath)
}
