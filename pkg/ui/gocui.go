package ui

import (
	"log"
	"sort"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type (
	KeybindingMap  map[string]*KeyBindingInfo
	KeybindingList []*KeyBindingInfo
	KeyBindingInfo struct {
		key    string
		action string
	}
)

type View struct {
	*gocui.View
	keybindingInfo KeybindingMap
}

type Gui struct {
	*gocui.Gui
}

var _gui *Gui

func (kb KeybindingMap) toList() KeybindingList {
	out := make(KeybindingList, 0)
	for _, kbm := range kb {
		out = append(out, kbm)
	}

	return out.Sorted()
}

func (t KeybindingList) Len() int { return len(t) }

func (t KeybindingList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t KeybindingList) Less(i, j int) bool {
	return strings.Compare(t[i].key, t[j].key) < 0
}

func (t KeybindingList) Sorted() KeybindingList {
	sort.Sort(t)
	return t
}

func newView(v *gocui.View) *View {
	return &View{
		View:           v,
		keybindingInfo: map[string]*KeyBindingInfo{},
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
	_gui.DeleteKeybindings(v.Name())
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

func (v *View) SetKeybindingInfo(key interface{}, description string) {
	keyStr := ""
	if s, ok := key.(rune); ok {
		keyStr = string(s)
	} else if s, ok := key.(gocui.Key); ok {
		keyStr = GetKeyStr(s)
	}

	v.keybindingInfo[keyStr] = &KeyBindingInfo{
		key:    keyStr,
		action: description,
	}
}

type KeyBindingBuilder struct {
	view *View
}

func NewKeybindingBuilder(view *View) *KeyBindingBuilder {
	return &KeyBindingBuilder{
		view: view,
	}
}

func (v *View) KeyBinding() *KeyBindingBuilder {
	return NewKeybindingBuilder(v)
}

func (kb *KeyBindingBuilder) set(key interface{}, description string, action func()) *KeyBindingBuilder {
	kb.setWithQuit(key, func() bool {
		action()
		return false
	}, description)
	return kb
}

func (kb *KeyBindingBuilder) setWithQuit(key interface{}, action func() bool, description string) *KeyBindingBuilder {
	name := ""

	if kb.view != nil {
		name = kb.view.Name()
		kb.view.SetKeybindingInfo(key, description)
	}

	if err := _gui.SetKeybinding(name, key, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
		if action() {
			return gocui.ErrQuit
		}
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	return kb
}

func GetKeyStr(key gocui.Key) string {
	for keyStr, k := range translate {
		if key == k {
			return keyStr
		}
	}
	return ""
}

// commenting some temporarily to avoid colliding keys
var translate = map[string]gocui.Key{
	"CtrlI":      gocui.KeyCtrlI,
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

	// "F1":             gocui.KeyF1,
	// "F2":             gocui.KeyF2,
	// "F3":             gocui.KeyF3,
	// "F4":             gocui.KeyF4,
	// "F5":             gocui.KeyF5,
	// "F6":             gocui.KeyF6,
	// "F7":             gocui.KeyF7,
	// "F8":             gocui.KeyF8,
	// "F9":             gocui.KeyF9,
	// "F10":            gocui.KeyF10,
	// "F11":            gocui.KeyF11,
	// "F12":            gocui.KeyF12,
	// "Insert":         gocui.KeyInsert,
	// "Delete":         gocui.KeyDelete,
	// "Home":           gocui.KeyHome,
	// "End":            gocui.KeyEnd,
	// "Pgup":           gocui.KeyPgup,
	// "Pgdn":           gocui.KeyPgdn,
	// "CtrlTilde":      gocui.KeyCtrlTilde,
	// "Ctrl2":          gocui.KeyCtrl2,
	// "CtrlSpace":      gocui.KeyCtrlSpace,
	// "CtrlA":          gocui.KeyCtrlA,
	// "CtrlB":          gocui.KeyCtrlB,
	// "CtrlC":          gocui.KeyCtrlC,
	// "CtrlD":          gocui.KeyCtrlD,
	// "CtrlE":          gocui.KeyCtrlE,
	// "CtrlF":          gocui.KeyCtrlF,
	// "CtrlG":          gocui.KeyCtrlG,
	// "Backspace":      gocui.KeyBackspace,
	// "Tab":            gocui.KeyTab,
	// "Backtab":        gocui.KeyBacktab,
	// "CtrlM":          gocui.KeyCtrlM,
	// "CtrlN":          gocui.KeyCtrlN,
	// "CtrlO":          gocui.KeyCtrlO,
	// "CtrlP":          gocui.KeyCtrlP,
	// "CtrlQ":          gocui.KeyCtrlQ,
	// "CtrlR":          gocui.KeyCtrlR,
	// "CtrlS":          gocui.KeyCtrlS,
	// "CtrlT":          gocui.KeyCtrlT,
	// "CtrlU":          gocui.KeyCtrlU,
	// "CtrlV":          gocui.KeyCtrlV,
	// "CtrlW":          gocui.KeyCtrlW,
	// "CtrlX":          gocui.KeyCtrlX,
	// "CtrlY":          gocui.KeyCtrlY,
	// "CtrlZ":          gocui.KeyCtrlZ,
	// "CtrlLsqBracket": gocui.KeyCtrlLsqBracket,
	// "Ctrl3":          gocui.KeyCtrl3,
	// "Ctrl4":          gocui.KeyCtrl4,
	// "CtrlBackslash":  gocui.KeyCtrlBackslash,
	// "Ctrl5":          gocui.KeyCtrl5,
	// "CtrlRsqBracket": gocui.KeyCtrlRsqBracket,
	// "Ctrl6":          gocui.KeyCtrl6,
	// "Ctrl7":          gocui.KeyCtrl7,
	// "CtrlSlash":      gocui.KeyCtrlSlash,
	// "CtrlUnderscore": gocui.KeyCtrlUnderscore,
	// "Space":          gocui.KeySpace,
	// "Backspace2":     gocui.KeyBackspace2,
	// "Ctrl8":          gocui.KeyCtrl8,
	// "Mouseleft":      gocui.MouseLeft,
	// "Mousemiddle":    gocui.MouseMiddle,
	// "Mouseright":     gocui.MouseRight,
	// "Mouserelease":   gocui.MouseRelease,
	// "MousewheelUp":   gocui.MouseWheelUp,
	// "MousewheelDown": gocui.MouseWheelDown,
}
