package utils

import (
	"errors"
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

	return DirExists("/Applications/Warp.app")
}

func IsItermInstalled() bool {
	if DetectOS() != Darwin {
		return false
	}

	return DirExists("/Applications/iTerm.app")
}

func OpenTerminal(path string) error {
	var cmd *exec.Cmd
	command := func(c ...string) *exec.Cmd {
		c = append(c, path)
		return exec.Command(c[0], c[1:]...)
	}
	switch DetectOS() {
	case Linux:
		cmd = command("xdg-open", "terminal")
	case Darwin:
		if IsWarpInstalled() {
			cmd = command("open", "-a", "warp")
		} else if IsItermInstalled() {
			cmd = command("open", "-a", "Iterm")
		} else {
			cmd = command("open", "-a", "Terminal")
		}
	default:
		return errors.New("unsupported OS")
	}

	if err := cmd.Start(); err != nil {
		return errors.New("failed to open terminal")
	}

	return nil
}
