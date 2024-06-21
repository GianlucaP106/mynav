package system

import "os"

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}
