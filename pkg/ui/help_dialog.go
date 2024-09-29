package ui

import (
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type helpDialog struct {
	view           *tui.View
	tableRenderer  *tui.TableRenderer[*tui.KeyBindingInfo]
	globalMappings []*tui.KeyBindingInfo
	mappings       []*tui.KeyBindingInfo
}

func newHelpViewEditor(up func(), down func(), enter func(), exit func()) gocui.EditorFunc {
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

func openHelpDialog(mappings []*tui.KeyBindingInfo, exit func()) *helpDialog {
	hv := &helpDialog{}
	hv.mappings = mappings

	x, _ := tui.ScreenSize()
	hv.view = tui.SetCenteredView(HelpDialog, (x*2)/3, 16, 0)
	hv.view.Editable = true
	hv.view.FrameColor = onFrameColor
	hv.view.Title = tui.WithSurroundingSpaces("Cheatsheet")
	ui.styleView(hv.view)

	prevView := tui.GetFocusedView()
	hv.view.Editor = newHelpViewEditor(func() {
		hv.tableRenderer.Up()
		hv.render()
	}, func() {
		hv.tableRenderer.Down()
		hv.render()
	}, func() {
	}, func() {
		hv.close()
		if prevView != nil {
			prevView.Focus()
		}
		exit()
	})

	sizeX, sizeY := hv.view.Size()
	hv.tableRenderer = tui.NewTableRenderer[*tui.KeyBindingInfo]()
	title := []string{
		"Key",
		"Action",
	}
	proportions := []float64{
		0.33,
		0.66,
	}

	hv.tableRenderer.InitTable(sizeX, sizeY, title, proportions)
	hv.refreshTable()
	hv.render()
	hv.view.Focus()
	return hv
}

func (hv *helpDialog) close() {
	hv.mappings = nil
	hv.view.Delete()
}

func (hv *helpDialog) refreshTable() {
	rows := make([][]string, 0)
	rowValues := make([]*tui.KeyBindingInfo, 0)
	for _, m := range hv.mappings {
		rowValues = append(rowValues, m)
		rows = append(rows, []string{
			m.Key,
			m.Action,
		})
	}

	for _, gm := range hv.globalMappings {
		rowValues = append(rowValues, gm)
		rows = append(rows, []string{
			gm.Key,
			gm.Action,
		})
	}

	hv.tableRenderer.FillTable(rows, rowValues)
}

func (hv *helpDialog) render() error {
	hv.view.Clear()
	hv.tableRenderer.Render(hv.view)
	return nil
}
