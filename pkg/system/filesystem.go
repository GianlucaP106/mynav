package system

import (
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

func SaveJson[T any](data *T, store string) error {
	if !Exists(store) {
		os.Create(store)
	}

	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(store, json, 0644)
}

func LoadJson[T any](store string) (*T, error) {
	file, err := os.Open(store)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		log.Panicln(err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data T
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func ShortenPath(path string, maxLength int) string {
	if len(path) <= maxLength {
		return path
	}

	dir, file := filepath.Split(path)
	dir = filepath.Clean(dir)

	ellipsis := "..."
	fileLen := len(file)
	dirLen := maxLength - fileLen - len(ellipsis)

	if dirLen <= 0 {
		return ellipsis + file[len(file)-maxLength+len(ellipsis):]
	}

	shortenedDir := dir[:dirLen] + ellipsis
	return filepath.Join(shortenedDir, file)
}
