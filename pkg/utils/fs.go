package utils

import (
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

func DirExists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}
