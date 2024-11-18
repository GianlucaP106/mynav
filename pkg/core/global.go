package core

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type GlobalConfigData struct {
	UpdateAsked time.Time `json:"update-asked"`
}

// GlobalConfig exposes crud on global configuration (~/.mynav)
type GlobalConfig struct {
	datasource *Datasource[GlobalConfigData]
}

func newGlobalConfig() (*GlobalConfig, error) {
	g := &GlobalConfig{}
	d, err := newDatasource(g.configPath(), &GlobalConfigData{})
	if err != nil {
		return nil, err
	}
	g.datasource = d

	return g, nil
}

func (gc *GlobalConfig) dirPath() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}

	return filepath.Join(dir, ".mynav")
}

func (gc *GlobalConfig) configPath() string {
	dir := gc.dirPath()
	return filepath.Join(dir, "config.json")
}

func (g *GlobalConfig) ConfigData() *GlobalConfigData {
	return g.datasource.Get()
}

func (gc *GlobalConfig) SetUpdateAsked(time time.Time) {
	data := gc.datasource.Get()
	data.UpdateAsked = time
	gc.datasource.Save(data)
}
