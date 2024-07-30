package core

import (
	"errors"
	"mynav/pkg/persistence"
	"mynav/pkg/system"
	"os"
	"path/filepath"
)

type WorkspaceDataSchema struct {
	Workspaces        map[string]*WorkspaceMetadata `json:"workspaces"`
	SelectedWorkspace string                        `json:"selected-workspace"`
}

type WorkspaceRepository struct {
	Container  *persistence.Container[Workspace]
	Datasource *persistence.Datasource[WorkspaceDataSchema]
}

func NewWorkspaceRepository(topics Topics, storePath string) *WorkspaceRepository {
	w := &WorkspaceRepository{}
	w.Datasource = persistence.NewDatasource[WorkspaceDataSchema](storePath)

	w.Datasource.LoadData()
	if w.Datasource.GetData() == nil {
		w.Datasource.SaveData(&WorkspaceDataSchema{
			Workspaces:        map[string]*WorkspaceMetadata{},
			SelectedWorkspace: "",
		})
	}

	w.LoadContainer(topics)
	return w
}

func (w *WorkspaceRepository) Save(workspace *Workspace) error {
	existing := w.Container.Get(workspace.ShortPath())
	if existing == nil {
		if err := system.CreateDir(workspace.Path); err != nil {
			return err
		}
	}

	w.Container.Set(workspace.ShortPath(), workspace)
	data := w.Datasource.GetData()
	data.Workspaces[workspace.ShortPath()] = workspace.Metadata
	w.Datasource.SaveData(data)
	return nil
}

func (w *WorkspaceRepository) Move(workspace *Workspace, newTopic *Topic) error {
	w.Container.Delete(workspace.ShortPath())
	w.DeleteMetadata(workspace)
	workspace.Topic = newTopic

	if err := w.Rename(workspace, workspace.Name); err != nil {
		return err
	}

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

	w.Container.Delete(workspace.ShortPath())
	w.DeleteMetadata(workspace)

	data := w.Datasource.GetData()
	if data.SelectedWorkspace == workspace.ShortPath() {
		data.SelectedWorkspace = newShortPath
		w.Datasource.SaveData(data)
	}

	workspace.Name = newName
	workspace.Path = newPath

	w.Container.Set(workspace.ShortPath(), workspace)
	w.SaveMetadata(workspace)

	return nil
}

func (w *WorkspaceRepository) Delete(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path); err != nil {
		return err
	}

	if wr := w.GetSelectedWorkspace(); wr != nil && wr.ShortPath() == workspace.ShortPath() {
		data := w.Datasource.GetData()
		data.SelectedWorkspace = ""
		w.Datasource.SaveData(data)
	}

	w.Container.Delete(workspace.ShortPath())
	w.DeleteMetadata(workspace)

	return nil
}

func (wr *WorkspaceRepository) Sync(w *persistence.Container[Workspace]) {
	data := wr.Datasource.GetData()
	for id, m := range data.Workspaces {
		if w.Get(id) == nil || m.Description == "" {
			delete(data.Workspaces, id)
		}
	}

	id := data.SelectedWorkspace
	if w.Get(id) == nil {
		data.SelectedWorkspace = ""
	}

	wr.Datasource.SaveData(data)
}

func (w *WorkspaceRepository) LoadContainer(topics Topics) {
	wc := persistence.NewContainer[Workspace]()
	w.Container = wc
	for _, topic := range topics {
		workspaceDirEntries := system.GetDirEntries(topic.Path)
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
			wc.Set(workspace.ShortPath(), workspace)
		}
	}
	w.Sync(wc)
}

func (w *WorkspaceRepository) Find(shortPath string) *Workspace {
	return w.Container.Get(shortPath)
}

func (w *WorkspaceRepository) FindByPath(path string) *Workspace {
	for _, w := range w.Container.All() {
		if w.Path == path {
			return w
		}
	}

	return nil
}

func (w *WorkspaceRepository) GetSelectedWorkspace() *Workspace {
	shortPath := w.Datasource.GetData().SelectedWorkspace
	return w.Find(shortPath)
}

func (wr *WorkspaceRepository) SetSelectedWorkspace(workspace *Workspace) error {
	data := wr.Datasource.GetData()
	data.SelectedWorkspace = workspace.ShortPath()
	if err := wr.Save(workspace); err != nil {
		return err
	}
	return nil
}

func (wr *WorkspaceRepository) SaveMetadata(w *Workspace) {
	data := wr.Datasource.GetData()
	data.Workspaces[w.ShortPath()] = w.Metadata
	wr.Datasource.SaveData(data)
}

func (wr *WorkspaceRepository) DeleteMetadata(w *Workspace) {
	data := wr.Datasource.GetData()
	delete(data.Workspaces, w.ShortPath())
	wr.Datasource.SaveData(data)
}

func (wr *WorkspaceRepository) GetMetadata(w *Workspace) *WorkspaceMetadata {
	return wr.Datasource.GetData().Workspaces[w.ShortPath()]
}
