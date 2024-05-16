package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type KeyBindingMappings struct {
	key    string
	action string
}

type HelpState struct {
	listRenderer   *ListRenderer
	viewName       string
	globalMappings []*KeyBindingMappings
	mappings       []*KeyBindingMappings
	active         bool
}

func (ui *UI) newHelpState(globalMappings []*KeyBindingMappings) *HelpState {
	return &HelpState{
		viewName:       "HelpView",
		globalMappings: globalMappings,
		listRenderer:   newListRenderer(0, 10, 0),
	}
}

func (ui *UI) initHelpView() *gocui.View {
	exists := false
	view := ui.getView(ui.help.viewName)
	exists = view != nil

	x, _ := ui.gui.Size()
	view = ui.setCenteredView(ui.help.viewName, x/2, 12, 0)

	if exists {
		return view
	}

	ui.keyBinding(ui.help.viewName).
		set(gocui.KeyEsc, func() {
			ui.closeHelpView()
		}).
		set('j', func() {
			ui.help.listRenderer.increment()
		}).
		set('k', func() {
			ui.help.listRenderer.decrement()
		}).
		set('?', func() {
			ui.closeHelpView()
		})

	return view
}

func (ui *UI) openHelpView(mappings []*KeyBindingMappings) {
	ui.help.mappings = mappings
	ui.refreshHelpListRenderer()
	ui.help.active = true
}

func (ui *UI) closeHelpView() {
	ui.help.active = false
	ui.help.mappings = nil
	ui.gui.DeleteView(ui.help.viewName)
}

func (ui *UI) refreshHelpListRenderer() {
	newSize := len(ui.help.mappings) + len(ui.help.globalMappings)
	if newSize != ui.help.listRenderer.listSize {
		ui.help.listRenderer.setListSize(newSize)
	}
}

func (ui *UI) formatHelpMessage(key *KeyBindingMappings, selected bool) string {
	view := ui.getView(ui.help.viewName)
	sizeX, _ := view.Size()

	color := func() color.Style {
		if selected {
			return color.New(color.Black, color.BgCyan)
		}
		return color.New(color.Blue)
	}()

	keyMap := withSpacePadding("[ "+key.key+" ]", sizeX/3)
	action := withSpacePadding(key.action, (sizeX*2)/3)
	return color.Sprint(keyMap + action)
}

func (ui *UI) renderHelpView() {
	if !ui.help.active {
		return
	}

	view := ui.initHelpView()

	mappings := append(ui.help.mappings, ui.help.globalMappings...)
	content := func() []string {
		out := make([]string, 0)
		ui.help.listRenderer.forEach(func(idx int) {
			helpMessage := mappings[idx]
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
