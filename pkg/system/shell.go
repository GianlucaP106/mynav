package system

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func OpenTerminalCmd(path string) (*exec.Cmd, error) {
	cmds := map[uint]func() string{
		Linux: func() string {
			return "xdg-open terminal"
		},
		Darwin: func() string {
			if IsItermInstalledMac() {
				return "open -a Iterm"
			} else if IsWarpInstalledMac() {
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
	return exec.Command(command[0], command[1:]...), nil
}

func CommandWithRedirect(command ...string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
