package system

import (
	"errors"
	"strings"
)

func GetNvimCmd(path string) []string {
	return []string{"nvim", path}
}

func GetNewTmuxSessionCmd(session string, path string) []string {
	return []string{"tmux", "new", "-s", session, "-c", path}
}

func GetAttachTmuxSessionCmd(session string) []string {
	return []string{"tmux", "a", "-t", session}
}

func GetKillTmuxSessionCmd(sessionName string) []string {
	return []string{"tmux", "kill-session", "-t", sessionName}
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

func GetUpdateSystemCmd() []string {
	return []string{"sh", "-c", "curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash"}
}
