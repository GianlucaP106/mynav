package ui

import (
	"github.com/awesome-gocui/gocui"
)

type Dialog interface {
	Name() string
	Render(ui *UI) error
}

func (ui *UI) InitDialogs() []gocui.Manager {
	ui.dialogs = map[string]Dialog{}
	ui.SetDialogs(
		newWorkspaceInfoDialogState(),
		newConfirmationDialogState(),
		newToastDialogState(),
		newEditorDialogState(),
		newHelpState(getKeyBindings("global")),
	)

	managers := []gocui.Manager{}
	for _, dialog := range ui.dialogs {
		manFunc := func(_ *gocui.Gui) error {
			return dialog.Render(ui)
		}
		managers = append(managers, gocui.ManagerFunc(manFunc))
	}

	return managers
}

func (ui *UI) SetDialog(d Dialog) {
	ui.dialogs[d.Name()] = d
}

func (ui *UI) SetDialogs(ds ...Dialog) {
	for _, d := range ds {
		ui.SetDialog(d)
	}
}

func GetDialog[T Dialog](ui *UI) T {
	for _, dialog := range ui.dialogs {
		v, ok := dialog.(T)
		if ok {
			return v
		}
	}

	panic("invalid dialog type")
}
