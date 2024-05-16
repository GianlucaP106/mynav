package utils

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
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

type TmuxSession struct {
	Name       string `json:"name"`
	NumWindows int    `json:"num-windows"`
}

func DeleteTmxSession(sessionName string) {
	err := exec.Command("tmux", "kill-session", "-t", sessionName).Run()
	if err != nil {
		log.Panicln(err)
	}
}

func GetTmuxSessions() map[string]*TmuxSession {
	stdout, err := exec.Command("tmux", "ls").Output()
	if err != nil {
		return map[string]*TmuxSession{}
	}

	out := map[string]*TmuxSession{}
	lines := strings.Split(string(stdout), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		name, numWindows := func() (string, int) {
			var name string
			var numWindows int

			nameSplit := strings.Split(line, ":")
			name = nameSplit[0]

			rest := strings.Join(nameSplit[1:], ":")
			numWindowSplit := strings.Split(rest, " ")

			var windowNumStr string
			for _, str := range numWindowSplit {
				if str != "" {
					windowNumStr = str
					break
				}
			}

			numWindows, err := strconv.Atoi(windowNumStr)
			if err != nil {
				log.Panicln(err)
			}

			return name, numWindows
		}()

		out[name] = &TmuxSession{
			Name:       name,
			NumWindows: numWindows,
		}
	}

	return out
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
