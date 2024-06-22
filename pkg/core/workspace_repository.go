package core

import (
	"errors"
	"mynav/pkg/filesystem"
	"os"
	"path/filepath"
)

type WorkspaceDataSchema struct {
	Workspaces        map[string]*WorkspaceMetadata `json:"workspaces"`
	SelectedWorkspace string                        `json:"selected-workspace"`
}

type WorkspaceRepository struct {
	Container  WorkspaceContainer
	Datasource *Datasource[WorkspaceDataSchema]
}

func NewWorkspaceRepository(topics Topics, storePath string) *WorkspaceRepository {
	w := &WorkspaceRepository{}
	w.Datasource = NewDatasource[WorkspaceDataSchema](storePath)

	w.Datasource.LoadData()
	if w.Datasource.Data == nil {
		w.Datasource.Data = &WorkspaceDataSchema{
			Workspaces:        map[string]*WorkspaceMetadata{},
			SelectedWorkspace: "",
		}
	}

	w.LoadContainer(topics)
	return w
}

func (w *WorkspaceRepository) Save(workspace *Workspace) error {
	existing := w.Container.Get(workspace.ShortPath())
	if existing == nil {
		if err := filesystem.CreateDir(workspace.Path); err != nil {
			return err
		}
	}

	w.Container.Set(workspace)
	w.Datasource.Data.Workspaces[workspace.ShortPath()] = workspace.Metadata
	w.Datasource.SaveData()

	return nil
}

func (w *WorkspaceRepository) Rename(workspace *Workspace, newName string) error {
	newShortPath := filepath.Join(workspace.Topic.Name, newName)
	if w.Container.Get(newShortPath) != nil {
		return errors.New("another workspace has this name already")
	}

	newPath := filepath.Join(workspace.Topic.Path, newName)
	if err := os.Rename(workspace.Path, newPath); err != nil {
		return err
	}

	w.Container.Delete(workspace)
	w.DeleteMetadata(workspace)

	if w.Datasource.Data.SelectedWorkspace == workspace.ShortPath() {
		w.Datasource.Data.SelectedWorkspace = newShortPath
	}

	workspace.Name = newName
	workspace.Path = newPath

	w.Container.Set(workspace)
	w.SaveMetadata(workspace)

	return nil
}

func (w *WorkspaceRepository) Delete(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path); err != nil {
		return err
	}

	if wr := w.GetSelectedWorkspace(); wr != nil && wr.ShortPath() == workspace.ShortPath() {
		w.Datasource.Data.SelectedWorkspace = ""
	}

	w.Container.Delete(workspace)
	w.DeleteMetadata(workspace)

	return nil
}

func (wr *WorkspaceRepository) Sync(w WorkspaceContainer) {
	for id, m := range wr.Datasource.Data.Workspaces {
		if w.Get(id) == nil || m.Description == "" {
			delete(wr.Datasource.Data.Workspaces, id)
		}
	}

	id := wr.Datasource.Data.SelectedWorkspace
	if w.Get(id) == nil {
		wr.Datasource.Data.SelectedWorkspace = ""
	}

	wr.Datasource.SaveData()
}

func (w *WorkspaceRepository) LoadContainer(topics Topics) {
	wc := NewWorkspaceContainer()
	w.Container = wc
	for _, topic := range topics {
		workspaceDirEntries := filesystem.GetDirEntries(topic.Path)
		for _, dirEntry := range workspaceDirEntries {
			if !dirEntry.IsDir() {
				continue
			}

			workspace := NewWorkspace(dirEntry.Name(), topic)
			metadata := w.GetMetadata(workspace)
			if metadata == nil {
				metadata = &WorkspaceMetadata{}
			}
			workspace.Metadata = metadata
			wc.Set(workspace)
		}
	}
	w.Sync(wc)
}

func (w *WorkspaceRepository) Find(shortPath string) *Workspace {
	return w.Container.Get(shortPath)
}

func (w *WorkspaceRepository) FindByPath(path string) *Workspace {
	for _, w := range w.Container {
		if w.Path == path {
			return w
		}
	}

	return nil
}

func (w *WorkspaceRepository) GetSelectedWorkspace() *Workspace {
	shortPath := w.Datasource.Data.SelectedWorkspace
	return w.Find(shortPath)
}

func (wr *WorkspaceRepository) SetSelectedWorkspace(workspace *Workspace) error {
	wr.Datasource.Data.SelectedWorkspace = workspace.ShortPath()
	if err := wr.Save(workspace); err != nil {
		return err
	}
	return nil
}

func (wr *WorkspaceRepository) SaveMetadata(w *Workspace) {
	wr.Datasource.Data.Workspaces[w.ShortPath()] = w.Metadata
	wr.Datasource.SaveData()
}

func (wr *WorkspaceRepository) DeleteMetadata(w *Workspace) {
	delete(wr.Datasource.Data.Workspaces, w.ShortPath())
	wr.Datasource.SaveData()
}

func (wr *WorkspaceRepository) GetMetadata(w *Workspace) *WorkspaceMetadata {
	return wr.Datasource.Data.Workspaces[w.ShortPath()]
}
