package ui

type Action struct {
	Command []string
	End     bool
}

func (ui *UI) setAction(command []string) {
	ui.action.Command = command
}

func (ui *UI) setActionEnd(command []string) {
	ui.action.End = true
	ui.action.Command = command
}

func (ui *UI) isActionReady() bool {
	return ui.action.Command != nil
}
