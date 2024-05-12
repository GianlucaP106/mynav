package utils

import (
	"errors"
	"runtime"
	"strings"
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

	return Exists("/Applications/Warp.app")
}

func IsItermInstalled() bool {
	if DetectOS() != Darwin {
		return false
	}

	return Exists("/Applications/iTerm.app")
}

func GetOpenTerminalCmd(path string) ([]string, error) {
	cmds := map[uint]func() string{
		Linux: func() string {
			return "xdg-open terminal"
		},
		Darwin: func() string {
			if IsItermInstalled() {
				return "open -a Iterm"
			} else if IsWarpInstalled() {
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

	return strings.Split(cmd(), " "), nil
}
