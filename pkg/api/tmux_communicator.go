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

func (tm *TmuxCommunicator) GetSessionPanes(ts *TmuxSession) []*TmuxPane {
	windowsRes, err := exec.Command("tmux", "list-windows", "-t", ts.Name).Output()
	if err != nil {
		log.Panicln(err)
	}

	var windows []string
	lines := strings.Split(string(windowsRes), "\n")
	for _, line := range lines {
		windowNumberSplit := strings.Split(line, ":")
		windowNumber := windowNumberSplit[0]
		windows = append(windows, windowNumber)
	}

	var panes []*TmuxPane
	for _, windowNum := range windows {
		paneRes, err := exec.Command("tmux", "list-panes", "-t", ts.Name+":"+windowNum, "-F", "#{pane_pid}:#{pane_id}").Output()
		if err != nil {
			log.Panicln(err)
		}

		paneLines := strings.Split(string(paneRes), "\n")
		for _, paneLine := range paneLines {
			pidNumberSplit := strings.Split(paneLine, ":%")
			if len(pidNumberSplit) != 2 {
				continue
			}

			pidStr := pidNumberSplit[0]
			paneNumberStr := pidNumberSplit[1]
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				log.Panicln(err)
			}

			paneNumber, err := strconv.Atoi(paneNumberStr)
			if err != nil {
				log.Panicln(err)
			}

			pane := NewTmuxPane(ts, pid, paneNumber)
			panes = append(panes, pane)
		}
	}

	return panes
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
