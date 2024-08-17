package system

import (
	"os"
	"os/exec"
	"runtime"

	"golang.design/x/clipboard"
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
	if err := clipboard.Init(); err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(s))
	return nil
}

func TimeFormat() string {
	return "Mon, 02 Jan 15:04:05"
}

func IsCurrentProcessHomeDir() bool {
	homeDir, _ := os.UserHomeDir()
	cwd, _ := os.Getwd()
	return homeDir == cwd
}
