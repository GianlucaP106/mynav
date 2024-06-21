package filesystem

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"time"
)

func GetDirEntries(d string) []fs.FileInfo {
	dir, err := os.Open(d)
	if err != nil {
		log.Panicln(err)
	}
	defer dir.Close()

	dirEntries, err := dir.Readdir(-1)
	if err != nil {
		log.Panicln(err)
	}
	return dirEntries
}

func GetLastModifiedTime(path string) (time.Time, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	return stat.ModTime(), nil
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateDir(path string) error {
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	return nil
}

func WriteFile(path string, b []byte) error {
	return os.WriteFile(path, b, 0644)
}

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
