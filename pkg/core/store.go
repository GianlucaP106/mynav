package core

import (
	"encoding/json"
	"io"
	"log"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

type Store struct {
	Some string `json:"some"`
}

func (fs *Filesystem) createStore() {
	os.Create(fs.GetDataPath())
}

func (fs *Filesystem) GetDataPath() string {
	return filepath.Join(fs.GetConfigPath(), "data.json")
}

func (fs *Filesystem) Save(store *Store) {
	dataPath := fs.GetDataPath()
	if !utils.Exists(dataPath) {
		fs.createStore()
	}

	json, err := json.Marshal(store)
	if err != nil {
		log.Panicln(err)
		return
	}

	if err := utils.WriteFile(dataPath, json); err != nil {
		log.Panicln(err)
		return
	}
}

func (fs *Filesystem) Load() *Store {
	file, err := os.Open(fs.GetDataPath())
	if err != nil {
		log.Panicln(err)
		return nil
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Panicln(err)
		return nil
	}

	var data Store
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Panicln(err)
		return nil
	}

	return &data
}
