package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type KeyBindingMapping struct {
	key    string
	action string
}

type HelpView struct {
	editor         Editor
	listRenderer   *ListRenderer
	globalMappings []*KeyBindingMapping
	mappings       []*KeyBindingMapping
}

var _ Dialog = &HelpView{}

const HelpStateName = "HelpView"

func newHelpState(globalMappings []*KeyBindingMapping) *HelpView {
	return &HelpView{
		globalMappings: globalMappings,
		listRenderer:   newListRenderer(0, 10, 0),
	}
}

func NewHelpViewEditor(up func(), down func(), enter func(), exit func()) gocui.EditorFunc {
	return gocui.EditorFunc(
		func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			switch {
			case key == gocui.KeyEnter:
				enter()
			case key == gocui.KeyEsc:
				exit()
			case ch == '?':
				exit()
			case ch == 'j':
				down()
			case ch == 'k':
				up()
			}
		})
}

func (hv *HelpView) Open(mappings []*KeyBindingMapping, exit func()) {
	hv.mappings = mappings
	x, _ := ScreenSize()
	prevView := GetFocusedView()
	hv.editor = NewHelpViewEditor(func() {
		hv.listRenderer.decrement()
	}, func() {
		hv.listRenderer.increment()
	}, func() {
	}, func() {
		hv.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		exit()
	})

	view := SetCenteredView(hv.Name(), x/2, 12, 0)
	view.Editable = true
	view.Editor = hv.editor
	view.FrameColor = gocui.ColorGreen

	FocusView(hv.Name())

	hv.refreshHelpListRenderer()
}

func (hv *HelpView) Close() {
	hv.mappings = nil
	DeleteView(hv.Name())
}

func (hv *HelpView) Name() string {
	return HelpStateName
}

func (hv *HelpView) refreshHelpListRenderer() {
	newSize := len(hv.mappings) + len(hv.globalMappings)
	if newSize != hv.listRenderer.listSize {
		hv.listRenderer.setListSize(newSize)
	}
}

func (hv *HelpView) formatHelpMessage(key *KeyBindingMapping, selected bool) string {
	view := GetInternalView(hv.Name())
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

func (hv *HelpView) Render(ui *UI) error {
	view := GetInternalView(hv.Name())
	if view == nil {
		return nil
	}

	mappings := append(hv.mappings, hv.globalMappings...)
	content := func() []string {
		out := make([]string, 0)
		hv.listRenderer.forEach(func(idx int) {
			helpMessage := mappings[idx]
			selected := idx == hv.listRenderer.selected
			out = append(out, hv.formatHelpMessage(helpMessage, selected))
		})
		return out
	}()

	view.Clear()
	sizeX, _ := view.Size()
	title := displayLine("Cheatsheet", Center, sizeX, color.New(color.White))
	fmt.Fprintln(view, title)
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
	return nil
}
