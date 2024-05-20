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

func (tm *TmuxCommunicator) GetSessions() TmuxSessions {
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

func (tm *TmuxCommunicator) DeleteSession(ts *TmuxSession) {
	err := exec.Command("tmux", "kill-session", "-t", ts.Name).Run()
	if err != nil {
		log.Panicln(err)
	}
}

func (tm *TmuxCommunicator) RenameSession(ts *TmuxSession, newName string) {
	if err := exec.Command("tmux", "rename-session", "-t", ts.Name, newName).Run(); err != nil {
		log.Panicln(err)
	}

	ts.Name = newName
}
