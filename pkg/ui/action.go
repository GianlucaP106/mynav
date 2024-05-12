package ui

type Action struct {
	Command []string
}

func (ui *UI) setAction(command []string) {
	if ui.action == nil {
		ui.action = &Action{}
	}
	ui.action.Command = command
}
