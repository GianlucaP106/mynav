package utils

import (
	"errors"
	"os"
	"strings"
)

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

	command := strings.Split(cmd(), " ")
	command = append(command, path)

	return command, nil
}

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func NvimCmd(path string) []string {
	return []string{"nvim", path}
}

func NewTmuxSessionCmd(session string, path string) []string {
	return []string{"tmux", "new", "-s", session, "-c", path}
}

func AttachTmuxSessionCmd(session string) []string {
	return []string{"tmux", "a", "-t", session}
}

func KillTmuxSessionCmd(sessionName string) []string {
	return []string{"tmux", "kill-session", "-t", sessionName}
}
