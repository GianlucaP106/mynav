package ui

type Action struct {
	Command []string
	End     bool
}

var action *Action

func SetAction(command []string) {
	action = &Action{
		Command: command,
	}
}

func SetActionEnd(command []string) {
	action = &Action{
		Command: command,
		End:     true,
	}
}

func IssActionReady() bool {
	return action != nil && action.Command != nil
}
