package utils

import "os"

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
