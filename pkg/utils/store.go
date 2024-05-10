package utils

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func Save[T any](data *T, store string) {
	if !Exists(store) {
		os.Create(store)
	}

	json, err := json.Marshal(data)
	if err != nil {
		log.Panicln(err)
		return
	}

	if err := WriteFile(store, json); err != nil {
		log.Panicln(err)
		return
	}
}

func Load[T any](store string) *T {
	file, err := os.Open(store)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Panicln(err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Panicln(err)
		return nil
	}

	var data T
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Panicln(err)
		return nil
	}

	return &data
}
