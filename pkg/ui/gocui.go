package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type View struct {
	*gocui.View
}

var _gui *gocui.Gui

func NewView(v *gocui.View) *View {
	return &View{
		View: v,
	}
}

func NewGui() *gocui.Gui {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	_gui = g
	return _gui
}

func GetView(name string) *View {
	view, err := _gui.View(name)
	if err != nil {
		return nil
	}
	return NewView(view)
}

func FocusViewInternal(name string) *View {
	v, err := _gui.SetCurrentView(name)
	if err != nil {
		return nil
	}
	return NewView(v)
}

func SetCenteredView(name string, sizeX int, sizeY int, verticalOffset int) *View {
	maxX, maxY := ScreenSize()
	view, _ := SetView(name, maxX/2-sizeX/2, maxY/2-sizeY/2+verticalOffset, maxX/2+sizeX/2, maxY/2+sizeY/2+verticalOffset, 0)
	return view
}

func DeleteView(name string) {
	_gui.DeleteView(name)
}

func ToggleCursor(c bool) {
	_gui.Cursor = c
}

func SetView(name string, x0 int, y0 int, x1 int, y1 int, overlaps byte) (*View, error) {
	v, err := _gui.SetView(name, x0, y0, x1, y1, overlaps)
	return NewView(v), err
}

func GetFocusedView() *View {
	v := _gui.CurrentView()
	if v != nil {
		return NewView(v)
	}

	return nil
}

func UpdateGui(f func(g *gocui.Gui) error) {
	_gui.Update(f)
}

func SetScreenManagers(managers ...gocui.Manager) {
	_gui.SetManager(managers...)
}

type KeyBindingBuilder struct {
	viewName string
}

func KeyBinding(viewName string) *KeyBindingBuilder {
	return &KeyBindingBuilder{
		viewName: viewName,
	}
}

func ScreenSize() (x int, y int) {
	return _gui.Size()
}

func (kb *KeyBindingBuilder) set(key interface{}, action func()) *KeyBindingBuilder {
	kb.setWithQuit(key, func() bool {
		action()
		return false
	})
	return kb
}

func (kb *KeyBindingBuilder) setWithQuit(key interface{}, action func() bool) *KeyBindingBuilder {
	if err := _gui.SetKeybinding(kb.viewName, key, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		if action() {
			return gocui.ErrQuit
		}
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	return kb
}
