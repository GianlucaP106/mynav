package core

import (
	"errors"
	"mynav/pkg/system"
	"os"
	"path/filepath"
)

type WorkspaceDataSchema struct {
	Workspaces        map[string]*WorkspaceMetadata `json:"workspaces"`
	SelectedWorkspace string                        `json:"selected-workspace"`
}

type WorkspaceRepository struct {
	container  *Container[Workspace]
	datasource *Datasource[WorkspaceDataSchema]
}

func NewWorkspaceRepository(topics Topics, storePath string) *WorkspaceRepository {
	w := &WorkspaceRepository{}
	ds, err := NewDatasource(storePath, &WorkspaceDataSchema{
		Workspaces:        map[string]*WorkspaceMetadata{},
		SelectedWorkspace: "",
	})
	if err != nil {
		return nil
	}

	w.datasource = ds
	w.LoadContainer(topics)
	return w
}

func (w *WorkspaceRepository) Save(workspace *Workspace) error {
	existing := w.container.Get(workspace.ShortPath())
	if existing == nil {
		if err := system.CreateDir(workspace.Path); err != nil {
			return err
		}
	}

	w.container.Set(workspace.ShortPath(), workspace)
	data := w.datasource.GetData()
	data.Workspaces[workspace.ShortPath()] = workspace.Metadata
	w.datasource.SaveData(data)
	return nil
}

func (w *WorkspaceRepository) Move(workspace *Workspace, newTopic *Topic) error {
	w.container.Delete(workspace.ShortPath())
	w.DeleteMetadata(workspace)
	workspace.Topic = newTopic

	if err := w.Rename(workspace, workspace.Name); err != nil {
		return err
	}

	return nil
}

func (w *WorkspaceRepository) Rename(workspace *Workspace, newName string) error {
	newShortPath := filepath.Join(workspace.Topic.Name, newName)
	if w.container.Get(newShortPath) != nil {
		return errors.New("another workspace has this name already")
	}

	newPath := filepath.Join(workspace.Topic.Path, newName)
	if err := os.Rename(workspace.Path, newPath); err != nil {
		return err
	}

	w.container.Delete(workspace.ShortPath())
	w.DeleteMetadata(workspace)

	data := w.datasource.GetData()
	if data.SelectedWorkspace == workspace.ShortPath() {
		data.SelectedWorkspace = newShortPath
		w.datasource.SaveData(data)
	}

	workspace.Name = newName
	workspace.Path = newPath

	w.container.Set(workspace.ShortPath(), workspace)
	w.SaveMetadata(workspace)

	return nil
}

func (w *WorkspaceRepository) Delete(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path); err != nil {
		return err
	}

	if wr := w.GetSelectedWorkspace(); wr != nil && wr.ShortPath() == workspace.ShortPath() {
		data := w.datasource.GetData()
		data.SelectedWorkspace = ""
		w.datasource.SaveData(data)
	}

	w.container.Delete(workspace.ShortPath())
	w.DeleteMetadata(workspace)

	return nil
}

func (wr *WorkspaceRepository) Sync(w *Container[Workspace]) {
	data := wr.datasource.GetData()
	for id, m := range data.Workspaces {
		if w.Get(id) == nil || m.Description == "" {
			delete(data.Workspaces, id)
		}
	}

	id := data.SelectedWorkspace
	if w.Get(id) == nil {
		data.SelectedWorkspace = ""
	}

	wr.datasource.SaveData(data)
}

func (w *WorkspaceRepository) LoadContainer(topics Topics) {
	wc := NewContainer[Workspace]()
	w.container = wc
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
	return w.container.Get(shortPath)
}

func (w *WorkspaceRepository) FindByPath(path string) *Workspace {
	for _, w := range w.container.All() {
		if w.Path == path {
			return w
		}
	}

	return nil
}

func (w *WorkspaceRepository) GetSelectedWorkspace() *Workspace {
	shortPath := w.datasource.GetData().SelectedWorkspace
	return w.Find(shortPath)
}

func (wr *WorkspaceRepository) SetSelectedWorkspace(workspace *Workspace) error {
	data := wr.datasource.GetData()
	data.SelectedWorkspace = workspace.ShortPath()
	if err := wr.Save(workspace); err != nil {
		return err
	}
	return nil
}

func (wr *WorkspaceRepository) SaveMetadata(w *Workspace) {
	data := wr.datasource.GetData()
	data.Workspaces[w.ShortPath()] = w.Metadata
	wr.datasource.SaveData(data)
}

func (wr *WorkspaceRepository) DeleteMetadata(w *Workspace) {
	data := wr.datasource.GetData()
	delete(data.Workspaces, w.ShortPath())
	wr.datasource.SaveData(data)
}

func (wr *WorkspaceRepository) GetMetadata(w *Workspace) *WorkspaceMetadata {
	return wr.datasource.GetData().Workspaces[w.ShortPath()]
}
