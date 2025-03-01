package tui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type View struct {
	*gocui.View
	Keybindings []*KeybindingInfo
}

type TUI struct {
	*gocui.Gui
	Keybindings []*KeybindingInfo
}

type ViewPosition struct {
	viewName string
	x0       int
	y0       int
	x1       int
	y1       int
	overlaps byte
}

func NewViewPosition(
	viewName string,
	x0 int,
	y0 int,
	x1 int,
	y1 int,
	overlaps byte,
) *ViewPosition {
	return &ViewPosition{
		viewName: viewName,
		x0:       x0,
		y0:       y0,
		x1:       x1,
		y1:       y1,
		overlaps: overlaps,
	}
}

func NewTui() *TUI {
	g, err := gocui.NewGui(gocui.OutputTrue, true)
	if err != nil {
		log.Panicln(err)
	}
	tui := &TUI{
		Gui:         g,
		Keybindings: make([]*KeybindingInfo, 0),
	}
	return tui
}

func newView(v *gocui.View) *View {
	return &View{
		View:        v,
		Keybindings: make([]*KeybindingInfo, 0),
	}
}

func Suspend() {
	gocui.Suspend()
}

func Resume() {
	gocui.Resume()
}

func (tui *TUI) Update(f func()) {
	tui.Gui.Update(func(g *gocui.Gui) error {
		f()
		return nil
	})
}

func (tui *TUI) SetManager(m func(t *TUI) error) {
	tui.Gui.SetManager(gocui.ManagerFunc(func(_ *gocui.Gui) error {
		return m(tui)
	}))
}

func (tui *TUI) SetCenteredView(name string, sizeX int, sizeY int, verticalOffset int) *View {
	maxX, maxY := tui.Size()
	p := NewViewPosition(name, maxX/2-sizeX/2, maxY/2-sizeY/2+verticalOffset, maxX/2+sizeX/2, maxY/2+sizeY/2+verticalOffset, 0)
	view := tui.SetView(p)
	return view
}

func (tui *TUI) SetView(p *ViewPosition) *View {
	v, _ := tui.Gui.SetView(p.viewName, p.x0, p.y0, p.x1, p.y1, p.overlaps)
	return newView(v)
}

func (tui *TUI) Resize(view *View, p *ViewPosition) *View {
	v, _ := tui.Gui.SetView(p.viewName, p.x0, p.y0, p.x1, p.y1, p.overlaps)
	view.View = v
	return view
}

func (tui *TUI) FocusedView() *View {
	v := tui.CurrentView()
	if v != nil {
		return newView(v)
	}

	return nil
}

func (tui *TUI) IsFocused(view *View) bool {
	v := tui.FocusedView()
	return v != nil && v.Name() == view.Name()
}

func (tui *TUI) FocusView(v *View) {
	tui.SetCurrentView(v.Name())
}

func (tui *TUI) DeleteView(v *View) {
	tui.Gui.DeleteView(v.Name())
	tui.DeleteKeybindings(v.Name())
}

type KeyBindingBuilder struct {
	view *View
	tui  *TUI
}

type KeybindingInfo struct {
	Key         string
	Description string
}

func newKeybindingInfo(key any, description string) *KeybindingInfo {
	keyStr := ""
	if s, ok := key.(rune); ok {
		keyStr = string(s)
	} else if s, ok := key.(gocui.Key); ok {
		keyStr = getKeyStr(s)
	}
	return &KeybindingInfo{
		Key:         keyStr,
		Description: description,
	}
}

func (tui *TUI) KeyBinding(v *View) *KeyBindingBuilder {
	return &KeyBindingBuilder{
		view: v,
		tui:  tui,
	}
}

func (kb *KeyBindingBuilder) Set(key interface{}, description string, action func()) *KeyBindingBuilder {
	kb.SetWithQuit(key, func() bool {
		action()
		return false
	}, description)
	return kb
}

func (kb *KeyBindingBuilder) SetWithQuit(key interface{}, action func() bool, description string) *KeyBindingBuilder {
	name := ""
	k := newKeybindingInfo(key, description)
	if kb.view != nil {
		name = kb.view.Name()
		kb.view.Keybindings = append(kb.view.Keybindings, k)
	} else {
		kb.tui.Keybindings = append(kb.tui.Keybindings, k)
	}

	if err := kb.tui.SetKeybinding(name, key, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		if action() {
			return gocui.ErrQuit
		}
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	return kb
}

func getKeyStr(key gocui.Key) string {
	for keyStr, k := range translate {
		if key == k {
			return keyStr
		}
	}
	return ""
}

var translate = map[string]gocui.Key{
	"CtrlJ":      gocui.KeyCtrlJ,
	"CtrlK":      gocui.KeyCtrlK,
	"CtrlL":      gocui.KeyCtrlL,
	"Enter":      gocui.KeyEnter,
	"ArrowUp":    gocui.KeyArrowUp,
	"ArrowDown":  gocui.KeyArrowDown,
	"ArrowLeft":  gocui.KeyArrowLeft,
	"ArrowRight": gocui.KeyArrowRight,
	"CtrlH":      gocui.KeyCtrlH,
	"Esc":        gocui.KeyEsc,
	"Tab":        gocui.KeyTab,
}
