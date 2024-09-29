package core

import (
	"mynav/pkg/system"
	"path/filepath"
)

type WorkspaceConfigurationData struct {
	Subworkspaces []string `json:"subworkspaces"`
}

type WorkspaceConfiguration struct {
	workspace  *Workspace
	datasource *Datasource[WorkspaceConfigurationData]
}

func NewWorkspaceConfiguration(w *Workspace) *WorkspaceConfiguration {
	wc := &WorkspaceConfiguration{}
	wc.workspace = w
	return wc
}

func (wc *WorkspaceConfiguration) init() error {
	configDirPath := wc.path()
	if !system.Exists(configDirPath) {
		if err := system.CreateDir(configDirPath); err != nil {
			return err
		}
	}

	workspaceConfigPath := filepath.Join(configDirPath, "config.json")
	d, err := NewDatasource[WorkspaceConfigurationData](workspaceConfigPath, nil)
	if err != nil {
		return err
	}
	wc.datasource = d

	return nil
}

func (wc *WorkspaceConfiguration) addSubworkspace(w *Workspace) error {
	data := wc.datasource.GetData()
	if data == nil {
		data = &WorkspaceConfigurationData{}
	}

	for _, v := range data.Subworkspaces {
		if v == w.Name {
			return nil
		}
	}

	data.Subworkspaces = append(data.Subworkspaces, w.Name)
	return wc.datasource.SaveData(data)
}

func (wc *WorkspaceConfiguration) removeSubworkspace(w *Workspace) {
	data := wc.datasource.GetData()
	for idx, v := range data.Subworkspaces {
		if v == w.Name {
			data.Subworkspaces = append(data.Subworkspaces[:idx], data.Subworkspaces[idx+1:]...)
		}
	}
	wc.datasource.SaveData(data)
	// TODO: delete thing if empty
}

func (wc *WorkspaceConfiguration) path() string {
	return filepath.Join(wc.workspace.Path, ".mynav")
}
