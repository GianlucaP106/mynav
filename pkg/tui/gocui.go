package tui

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
		Key    string
		Action string
	}
)

type View struct {
	*gocui.View
	keybindingInfo KeybindingMap
}

type Tui struct {
	*gocui.Gui
	globalKeybindingInfo KeybindingMap
}

var _tui *Tui

func (kb KeybindingMap) ToList() KeybindingList {
	out := make(KeybindingList, 0)
	for _, kbm := range kb {
		out = append(out, kbm)
	}

	return out.Sorted()
}

func (kb KeybindingMap) Set(key interface{}, description string) {
	keyStr := ""
	if s, ok := key.(rune); ok {
		keyStr = string(s)
	} else if s, ok := key.(gocui.Key); ok {
		keyStr = GetKeyStr(s)
	}

	if keyStr != "" {
		kb[keyStr] = &KeyBindingInfo{
			Key:    keyStr,
			Action: description,
		}
	}
}

func (t KeybindingList) Len() int { return len(t) }

func (t KeybindingList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t KeybindingList) Less(i, j int) bool {
	return strings.Compare(t[i].Key, t[j].Key) < 0
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

func NewTui() *Tui {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	_tui = &Tui{
		Gui:                  g,
		globalKeybindingInfo: make(KeybindingMap),
	}
	return _tui
}

func ScreenSize() (x int, y int) {
	return _tui.Size()
}

func UpdateTui(f func(g *Tui) error) {
	_tui.Update(func(g *gocui.Gui) error {
		return f(&Tui{
			Gui: g,
		})
	})
}

func RunAction(action func()) {
	gocui.Suspend()
	action()
	gocui.Resume()
}

func SetManagerFunctions(managers ...gocui.Manager) {
	_tui.SetManager(managers...)
}

func ToggleCursor(c bool) {
	_tui.Cursor = c
}

func GetView(name string) *View {
	view, err := _tui.View(name)
	if err != nil {
		return nil
	}
	return newView(view)
}

func SetFocusView(name string) *View {
	v, err := _tui.SetCurrentView(name)
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
	v, _ := _tui.SetView(name, x0, y0, x1, y1, overlaps)
	return newView(v)
}

func GetFocusedView() *View {
	v := _tui.CurrentView()
	if v != nil {
		return newView(v)
	}

	return nil
}

func (v *View) GetKeybindings() KeybindingList {
	viewKeys := v.keybindingInfo.ToList()
	globalKeys := _tui.globalKeybindingInfo.ToList()
	viewKeys = append(viewKeys, globalKeys...)
	return viewKeys
}

func (v *View) Delete() {
	_tui.DeleteView(v.Name())
	_tui.DeleteKeybindings(v.Name())
}

func (vw *View) Focus() *View {
	v, err := _tui.SetCurrentView(vw.Name())
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
	_tui.SetViewOnBottom(v.Name())
}

func (v *View) SendToFront() {
	_tui.SetViewOnTop(v.Name())
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

func (kb *KeyBindingBuilder) Set(key interface{}, description string, action func()) *KeyBindingBuilder {
	kb.SetWithQuit(key, func() bool {
		action()
		return false
	}, description)
	return kb
}

func (kb *KeyBindingBuilder) SetWithQuit(key interface{}, action func() bool, description string) *KeyBindingBuilder {
	name := ""

	if kb.view != nil {
		name = kb.view.Name()
		kb.view.keybindingInfo.Set(key, description)
	} else {
		_tui.globalKeybindingInfo.Set(key, description)
	}

	if err := _tui.SetKeybinding(name, key, gocui.ModNone, func(_ *gocui.Gui, _ *gocui.View) error {
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
