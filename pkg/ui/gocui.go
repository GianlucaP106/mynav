package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

var gui *gocui.Gui

func NewGui() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	gui = g
	return gui
}

func GetInternalView(name string) *gocui.View {
	view, err := gui.View(name)
	if err != nil {
		return nil
	}
	return view
}

func FocusView(name string) *gocui.View {
	v, err := gui.SetCurrentView(name)
	if err != nil {
		return nil
	}
	return v
}

func SetCenteredView(name string, sizeX int, sizeY int, verticalOffset int) *gocui.View {
	maxX, maxY := ScreenSize()
	view, _ := SetView(name, maxX/2-sizeX/2, maxY/2-sizeY/2+verticalOffset, maxX/2+sizeX/2, maxY/2+sizeY/2+verticalOffset, 0)
	return view
}

func DeleteView(name string) {
	gui.DeleteView(name)
}

func ToggleCursor(c bool) {
	gui.Cursor = c
}

func SetView(name string, x0 int, y0 int, x1 int, y1 int, overlaps byte) (*gocui.View, error) {
	return gui.SetView(name, x0, y0, x1, y1, overlaps)
}

func GetFocusedView() *gocui.View {
	return gui.CurrentView()
}

func UpdateGui(f func(g *gocui.Gui) error) {
	gui.Update(f)
}

func SetScreenManagers(managers ...gocui.Manager) {
	gui.SetManager(managers...)
}

type KeyBindingBuilder struct {
	viewName string
}

func KeyBinding(viewName string) *KeyBindingBuilder {
	return &KeyBindingBuilder{
		viewName: viewName,
	}
}

func (kb *KeyBindingBuilder) setKeybinding(
	viewName string,
	key interface{},
	handler func(g *gocui.Gui, v *gocui.View) error,
) *KeyBindingBuilder {
	if err := gui.SetKeybinding(viewName, key, gocui.ModNone, handler); err != nil {
		log.Panicln(err)
	}
	return kb
}

func (kb *KeyBindingBuilder) set(key interface{}, action func()) *KeyBindingBuilder {
	kb.setKeybinding(kb.viewName, key, func(g *gocui.Gui, v *gocui.View) error {
		action()
		return nil
	})
	return kb
}

func ScreenSize() (x int, y int) {
	return gui.Size()
}
