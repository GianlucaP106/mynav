package api

import (
	"mynav/pkg/utils"
	"os"
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
		if err := utils.CreateDir(workspace.Path); err != nil {
			return err
		}
	}

	w.WorkspaceContainer.Set(workspace)
	w.WorkspaceDatasource.SaveMetadata(workspace)
	w.WorkspaceDatasource.SetSelectedWorkspace(workspace)

	return nil
}

func (w *WorkspaceRepository) LoadContainer(topics Topics) {
	wc := NewWorkspaceContainer()
	w.WorkspaceContainer = wc
	for _, topic := range topics {
		workspaceDirEntries := utils.GetDirEntries(topic.Path)
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
	shortPath := w.WorkspaceDatasource.WorkspaceStoreSchema.SelectedWorkspace
	return w.Find(shortPath)
}
