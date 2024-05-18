package api

import (
	"mynav/pkg/utils"
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
	utils.Save(c.Data, c.StorePath)
}

func (c *ConfigurationDatasource) LoadStore() {
	c.Data = utils.Load[ConfigurationDataSchema](c.StorePath)
	if c.Data == nil {
		c.Data = &ConfigurationDataSchema{
			UpdateAsked: nil,
		}
		c.SaveStore()
	}
}
