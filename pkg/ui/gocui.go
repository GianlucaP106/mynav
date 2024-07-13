package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type View struct {
	*gocui.View
}

type Gui struct {
	*gocui.Gui
}

var _gui *Gui

func newView(v *gocui.View) *View {
	return &View{
		View: v,
	}
}

func NewGui() *Gui {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	_gui = &Gui{
		Gui: g,
	}
	return _gui
}

func ScreenSize() (x int, y int) {
	return _gui.Size()
}

func UpdateGui(f func(g *Gui) error) {
	_gui.Update(func(g *gocui.Gui) error {
		return f(&Gui{
			Gui: g,
		})
	})
}

func SetManagerFunctions(managers ...gocui.Manager) {
	_gui.SetManager(managers...)
}

func ToggleCursor(c bool) {
	_gui.Cursor = c
}

func GetView(name string) *View {
	view, err := _gui.View(name)
	if err != nil {
		return nil
	}
	return newView(view)
}

func SetFocusView(name string) *View {
	v, err := _gui.SetCurrentView(name)
	if err != nil {
		return nil
	}
	return newView(v)
}

func SetCenteredView(name string, sizeX int, sizeY int, verticalOffset int) *View {
	maxX, maxY := ScreenSize()
	view := SetView(name, maxX/2-sizeX/2, maxY/2-sizeY/2+verticalOffset, maxX/2+sizeX/2, maxY/2+sizeY/2+verticalOffset, 0)
	return view
}

func SetView(name string, x0 int, y0 int, x1 int, y1 int, overlaps byte) *View {
	v, _ := _gui.SetView(name, x0, y0, x1, y1, overlaps)
	return newView(v)
}

func GetFocusedView() *View {
	v := _gui.CurrentView()
	if v != nil {
		return newView(v)
	}

	return nil
}

func (v *View) Delete() {
	_gui.DeleteView(v.Name())
}

func (vw *View) Focus() *View {
	v, err := _gui.SetCurrentView(vw.Name())
	if err != nil {
		return nil
	}
	return newView(v)
}

func (vw *View) IsFocused() bool {
	v := GetFocusedView()
	return v != nil && v.Name() == vw.Name()
}

func (v *View) SendToBack() {
	_gui.SetViewOnBottom(v.Name())
}

func (v *View) SendToFront() {
	_gui.SetViewOnTop(v.Name())
}

type KeyBindingBuilder struct {
	name string
}

func NewKeybindingBuilder(name string) *KeyBindingBuilder {
	return &KeyBindingBuilder{
		name: name,
	}
}

func (v *View) KeyBinding() *KeyBindingBuilder {
	return NewKeybindingBuilder(v.Name())
}

func (kb *KeyBindingBuilder) set(key interface{}, action func()) *KeyBindingBuilder {
	kb.setWithQuit(key, func() bool {
		action()
		return false
	})
	return kb
}

func (kb *KeyBindingBuilder) setWithQuit(key interface{}, action func() bool) *KeyBindingBuilder {
	if err := _gui.SetKeybinding(kb.name, key, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		if action() {
			return gocui.ErrQuit
		}
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	return kb
}
