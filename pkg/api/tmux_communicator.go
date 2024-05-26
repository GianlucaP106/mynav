package api

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type TmuxCommunicator struct{}

func NewTmuxCommunicator() *TmuxCommunicator {
	return &TmuxCommunicator{}
}

func (tm *TmuxCommunicator) GetSessions() map[string]*TmuxSession {
	out := map[string]*TmuxSession{}

	stdout, err := exec.Command("tmux", "ls").Output()
	if err != nil {
		return out
	}

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

func (tm *TmuxCommunicator) KillSession(name string) error {
	err := exec.Command("tmux", "kill-session", "-t", name).Run()
	if err != nil {
		return err
	}

	return nil
}

func (tm *TmuxCommunicator) RenameSession(oldName string, newName string) error {
	if err := exec.Command("tmux", "rename-session", "-t", oldName, newName).Run(); err != nil {
		return err
	}

	return nil
}
