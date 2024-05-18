package ui

type Action struct {
	Command []string
	End     bool
}

func (ui *UI) setAction(command []string) {
	if ui.action == nil {
		ui.action = &Action{}
	}
	ui.action.Command = command
}

func (ui *UI) setActionEnd(command []string) {
	if ui.action == nil {
		ui.action = &Action{}
	}
	ui.action.End = true
	ui.action.Command = command
}
