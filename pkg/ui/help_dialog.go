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
	tableRenderer  *TableRenderer
	globalMappings []*KeyBindingMapping
	mappings       []*KeyBindingMapping
}

var _ Dialog = &HelpView{}

const HelpStateName = "HelpView"

func newHelpState(globalMappings []*KeyBindingMapping) *HelpView {
	return &HelpView{
		globalMappings: globalMappings,
		tableRenderer:  NewTableRenderer(),
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
		hv.tableRenderer.Up()
	}, func() {
		hv.tableRenderer.Down()
	}, func() {
	}, func() {
		hv.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		exit()
	})

	view := SetCenteredView(hv.Name(), (x*2)/3, 16, 0)
	view.Editable = true
	view.Editor = hv.editor
	view.FrameColor = gocui.ColorGreen

	sizeX, _ := view.Size()
	title := []string{
		"Key",
		"Action",
	}
	proportions := []float64{
		0.33,
		0.66,
	}

	hv.tableRenderer.InitTable(sizeX, 13, title, proportions)
	hv.refreshTable()
	FocusView(hv.Name())
}

func (hv *HelpView) Close() {
	hv.mappings = nil
	DeleteView(hv.Name())
}

func (hv *HelpView) Name() string {
	return HelpStateName
}

func (hv *HelpView) refreshTable() {
	rows := make([][]string, 0)
	for _, m := range hv.mappings {
		rows = append(rows, []string{
			m.key,
			m.action,
		})
	}

	for _, gm := range hv.globalMappings {
		rows = append(rows, []string{
			gm.key,
			gm.action,
		})
	}
	hv.tableRenderer.FillTable(rows)
}

func (hv *HelpView) Render(ui *UI) error {
	view := GetInternalView(hv.Name())
	if view == nil {
		return nil
	}

	view.Clear()
	sizeX, _ := view.Size()
	title := displayLine("Cheatsheet", Center, sizeX, color.New(color.White))
	fmt.Fprintln(view, title)
	hv.tableRenderer.Render(view)
	return nil
}
