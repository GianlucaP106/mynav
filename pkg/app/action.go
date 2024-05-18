package app

import (
	"fmt"
	"mynav/pkg/ui"
	"os"
	"os/exec"
)

func handleAction(action *ui.Action) {
	if action.Command != nil {
		cmd := exec.Command(action.Command[0], action.Command[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
		if action.End {
			fmt.Println("Succesfully updated mynav")
			os.Exit(0)
		}
	}
}
