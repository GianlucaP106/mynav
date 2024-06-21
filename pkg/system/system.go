package system

import (
	"mynav/pkg/filesystem"
	"os/exec"
	"runtime"
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

func IsWarpInstalled() bool {
	if DetectOS() != Darwin {
		return false
	}

	return filesystem.Exists("/Applications/Warp.app")
}

func IsItermInstalled() bool {
	if DetectOS() != Darwin {
		return false
	}

	return filesystem.Exists("/Applications/iTerm.app")
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
