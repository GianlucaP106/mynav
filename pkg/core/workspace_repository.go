package core

import (
	"errors"
	"mynav/pkg/filesystem"
	"os"
	"path/filepath"
)

type WorkspaceRepository struct {
	WorkspaceContainer  WorkspaceContainer
	WorkspaceDatasource *WorkspaceDatasource
}

func NewWorkspaceRepository(topics Topics, storePath string) *WorkspaceRepository {
	w := &WorkspaceRepository{}
	w.WorkspaceDatasource = NewWorkspaceDatasource(storePath)
	w.LoadContainer(topics)
	return w
}

func (w *WorkspaceRepository) Save(workspace *Workspace) error {
	existing := w.WorkspaceContainer.Get(workspace.ShortPath())
	if existing == nil {
		if err := filesystem.CreateDir(workspace.Path); err != nil {
			return err
		}
	}

	w.WorkspaceContainer.Set(workspace)
	w.WorkspaceDatasource.Data.Workspaces[workspace.ShortPath()] = workspace.Metadata
	w.WorkspaceDatasource.SaveStore()

	return nil
}

func (w *WorkspaceRepository) Rename(workspace *Workspace, newName string) error {
	newShortPath := filepath.Join(workspace.Topic.Name, newName)
	if w.WorkspaceContainer.Get(newShortPath) != nil {
		return errors.New("another workspace has this name already")
	}

	newPath := filepath.Join(workspace.Topic.Path, newName)
	if err := os.Rename(workspace.Path, newPath); err != nil {
		return err
	}

	w.WorkspaceContainer.Delete(workspace)
	w.WorkspaceDatasource.DeleteMetadata(workspace)

	if w.WorkspaceDatasource.Data.SelectedWorkspace == workspace.ShortPath() {
		w.WorkspaceDatasource.Data.SelectedWorkspace = newShortPath
	}

	workspace.Name = newName
	workspace.Path = newPath

	w.WorkspaceContainer.Set(workspace)
	w.WorkspaceDatasource.SaveMetadata(workspace)

	return nil
}

func (w *WorkspaceRepository) LoadContainer(topics Topics) {
	wc := NewWorkspaceContainer()
	w.WorkspaceContainer = wc
	for _, topic := range topics {
		workspaceDirEntries := filesystem.GetDirEntries(topic.Path)
		for _, dirEntry := range workspaceDirEntries {
			if !dirEntry.IsDir() {
				continue
			}

			workspace := NewWorkspace(dirEntry.Name(), topic)
			metadata := w.WorkspaceDatasource.GetMetadata(workspace)
			if metadata == nil {
				metadata = &WorkspaceMetadata{}
			}
			workspace.Metadata = metadata
			wc.Set(workspace)
		}
	}
	w.WorkspaceDatasource.Sync(wc)
}

func (w *WorkspaceRepository) Delete(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path); err != nil {
		return err
	}

	if wr := w.GetSelectedWorkspace(); wr != nil && wr.ShortPath() == workspace.ShortPath() {
		w.WorkspaceDatasource.Data.SelectedWorkspace = ""
	}

	w.WorkspaceContainer.Delete(workspace)
	w.WorkspaceDatasource.DeleteMetadata(workspace)

	return nil
}

func (w *WorkspaceRepository) GetContainer() WorkspaceContainer {
	return w.WorkspaceContainer
}

func (w *WorkspaceRepository) Find(shortPath string) *Workspace {
	return w.WorkspaceContainer.Get(shortPath)
}

func (w *WorkspaceRepository) FindByPath(path string) *Workspace {
	for _, w := range w.WorkspaceContainer {
		if w.Path == path {
			return w
		}
	}

	return nil
}

func (w *WorkspaceRepository) GetSelectedWorkspace() *Workspace {
	shortPath := w.WorkspaceDatasource.Data.SelectedWorkspace
	return w.Find(shortPath)
}

func (w *WorkspaceRepository) SetSelectedWorkspace(workspace *Workspace) error {
	w.WorkspaceDatasource.SetSelectedWorkspace(workspace)
	if err := w.Save(workspace); err != nil {
		return err
	}
	return nil
}
