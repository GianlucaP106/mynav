package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type KeyBindingMessage struct {
	key    string
	action string
}

type HelpState struct {
	listRenderer *ListRenderer
	viewName     string
	messages     []*KeyBindingMessage
	active       bool
}

func newHelpState() *HelpState {
	return &HelpState{
		viewName: "HelpView",
	}
}

func (ui *UI) initHelpView() *gocui.View {
	x, _ := ui.gui.Size()
	view := ui.getView(ui.help.viewName)
	if view != nil {
		return view
	}
	ui.help.messages = []*KeyBindingMessage{
		{
			key:    "q",
			action: "Quit",
		},
		{
			key:    "j",
			action: "Move down",
		},
		{
			key:    "k",
			action: "Move up",
		},
		{
			key:    "a",
			action: "Create entry",
		},
		{
			key:    "d",
			action: "Delete entry",
		},
		{
			key:    "enter",
			action: "Select",
		},
		{
			key:    "esc",
			action: "Go back",
		},
	}

	numKeymaps := len(ui.help.messages)
	view = ui.setCenteredView(ui.help.viewName, x/2, numKeymaps+5, 0)
	_, sizeY := view.Size()
	ui.help.listRenderer = newListRenderer(0, sizeY, numKeymaps)

	ui.keyBinding(ui.help.viewName).
		set(gocui.KeyEsc, func() {
			ui.closeHelpView()
		}).
		set('j', func() {
			ui.help.listRenderer.increment()
		}).
		set('k', func() {
			ui.help.listRenderer.decrement()
		})
	return view
}

func (ui *UI) openHelpView() {
	ui.help.active = true
}

func (ui *UI) closeHelpView() {
	ui.help.active = false
	ui.gui.DeleteView(ui.help.viewName)
}

func (ui *UI) formatHelpMessage(key *KeyBindingMessage, selected bool) string {
	view := ui.getView(ui.help.viewName)
	sizeX, _ := view.Size()

	color := func() color.Style {
		if selected {
			return color.New(color.Black, color.BgCyan)
		}
		return color.New(color.Blue)
	}()

	keyMap := withSpacePadding(key.key, sizeX/3)
	action := withSpacePadding(key.action, (sizeX*2)/3)
	return color.Sprint(keyMap + action)
}

func (ui *UI) renderHelpView() {
	if !ui.help.active {
		return
	}
	view := ui.initHelpView()

	content := func() []string {
		out := make([]string, 0)
		ui.help.listRenderer.forEach(func(idx int) {
			helpMessage := ui.help.messages[idx]
			selected := idx == ui.help.listRenderer.selected
			out = append(out, ui.formatHelpMessage(helpMessage, selected))
		})
		return out
	}()

	ui.focusView(ui.help.viewName)
	view.Clear()
	sizeX, _ := view.Size()
	title := displayLine("Cheatsheet", Center, sizeX, color.New(color.White))
	fmt.Fprintln(view, title)
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
}
