package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

func (ui *UI) getView(name string) *gocui.View {
	view, err := ui.gui.View(name)
	if err != nil {
		return nil
	}
	return view
}

func (ui *UI) focusView(name string) *gocui.View {
	v, err := ui.gui.SetCurrentView(name)
	if err != nil {
		return nil
	}
	return v
}

func (ui *UI) toggleCursor(on bool) {
	ui.gui.Cursor = on
}

type KeyBinding struct {
	ui       *UI
	viewName string
}

func (ui *UI) keyBinding(viewName string) *KeyBinding {
	return &KeyBinding{
		viewName: viewName,
		ui:       ui,
	}
}

func (kb *KeyBinding) setKeybinding(
	viewName string,
	key interface{},
	handler func(g *gocui.Gui, v *gocui.View) error,
) *KeyBinding {
	if err := kb.ui.gui.SetKeybinding(viewName, key, gocui.ModNone, handler); err != nil {
		log.Panicln(err)
	}
	return kb
}

func (kb *KeyBinding) set(key interface{}, action func()) *KeyBinding {
	kb.setKeybinding(kb.viewName, key, func(g *gocui.Gui, v *gocui.View) error {
		action()
		return nil
	})
	return kb
}
