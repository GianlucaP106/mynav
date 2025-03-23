package core

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
)

type OS = uint

const (
	Darwin OS = iota
	Linux
	Unsuported
)

func DetectOS() OS {
	switch runtime.GOOS {
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	default:
		return Unsuported
	}
}

func IsWarpInstalledMac() bool {
	if DetectOS() != Darwin {
		return false
	}

	return Exists("/Applications/Warp.app")
}

func IsItermInstalledMac() bool {
	if DetectOS() != Darwin {
		return false
	}

	return Exists("/Applications/iTerm.app")
}

func OpenBrowser(url string) error {
	var cmd string

	switch DetectOS() {
	case Linux:
		cmd = "xdg-open"
	case Darwin:
		cmd = "open"
	}

	return exec.Command(cmd, url).Start()
}

func CopyToClip(s string) error {
	return clipboard.WriteAll(s)
}

func DoesProgramExist(program string) bool {
	_, err := exec.LookPath(program)
	return err == nil
}

func OpenLazygit(path string) error {
	return exec.Command("lazygit", "-p", path).Run()
}

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

func OpenTerminalCmd(path string) (*exec.Cmd, error) {
	cmds := map[uint]func() string{
		Linux: func() string {
			return "xdg-open terminal"
		},
		Darwin: func() string {
			if IsItermInstalledMac() {
				return "open -a Iterm"
			} else if IsWarpInstalledMac() {
				return "open -a warp"
			} else {
				return "open -a Terminal"
			}
		},
	}

	os := DetectOS()
	cmd := cmds[os]
	if cmd == nil {
		return nil, errors.New("os not currently supported")
	}

	command := strings.Split(cmd(), " ")
	command = append(command, path)
	return exec.Command(command[0], command[1:]...), nil
}

func CommandWithRedirect(command ...string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
