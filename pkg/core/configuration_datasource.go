package core

import (
	"mynav/pkg/filesystem"
	"time"
)

type ConfigurationDataSchema struct {
	UpdateAsked *time.Time `json:"update-asked"`
}

type ConfigurationDatasource struct {
	Data      *ConfigurationDataSchema
	StorePath string
}

func NewConfigurationDatasource(storePath string) *ConfigurationDatasource {
	c := &ConfigurationDatasource{
		StorePath: storePath,
	}

	c.LoadStore()
	return c
}

func (c *ConfigurationDatasource) SaveStore() {
	filesystem.Save(c.Data, c.StorePath)
}

func (c *ConfigurationDatasource) LoadStore() {
	c.Data = filesystem.Load[ConfigurationDataSchema](c.StorePath)
	if c.Data == nil {
		c.Data = &ConfigurationDataSchema{
			UpdateAsked: nil,
		}
	}
}
