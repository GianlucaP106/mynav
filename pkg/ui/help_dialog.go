package ui

import (
	"fmt"
	"mynav/pkg/constants"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type HelpDialog struct {
	view           *View
	tableRenderer  *TableRenderer
	globalMappings []*KeyBindingInfo
	mappings       []*KeyBindingInfo
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

func OpenHelpView(mappings []*KeyBindingInfo, exit func()) *HelpDialog {
	hv := &HelpDialog{}
	hv.mappings = mappings

	x, _ := ScreenSize()
	hv.view = SetCenteredView(constants.HelpDialogName, (x*2)/3, 16, 0)
	hv.view.Editable = true
	hv.view.FrameColor = gocui.ColorGreen

	prevView := GetFocusedView()
	hv.view.Editor = NewHelpViewEditor(func() {
		hv.tableRenderer.Up()
		hv.render()
	}, func() {
		hv.tableRenderer.Down()
		hv.render()
	}, func() {
	}, func() {
		hv.Close()
		if prevView != nil {
			prevView.Focus()
		}
		exit()
	})

	sizeX, _ := hv.view.Size()
	hv.tableRenderer = NewTableRenderer()
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
	hv.render()
	hv.view.Focus()
	return hv
}

func (hv *HelpDialog) Close() {
	hv.mappings = nil
	hv.view.Delete()
}

func (hv *HelpDialog) refreshTable() {
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

func (hv *HelpDialog) render() error {
	hv.view.Clear()
	sizeX, _ := hv.view.Size()
	title := displayLine("Cheatsheet", Center, sizeX, color.New(color.White))
	fmt.Fprintln(hv.view, title)
	hv.tableRenderer.Render(hv.view)
	return nil
}
