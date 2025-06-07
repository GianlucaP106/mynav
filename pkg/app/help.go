package app

import (
	"fmt"
	"sort"

	"github.com/GianlucaP106/mynav/pkg/tui"
	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type Help struct {
	view *tui.View

	// tables for table and local keys
	table *tui.TableRenderer[*tui.KeybindingInfo]
}

func help(v *tui.View) {
	h := &Help{}
	h.view = a.ui.SetCenteredView(HelpDialog, 80, 20, 0, 0)
	h.view.Title = fmt.Sprintf(" %s ", "Key bindings")
	a.styleView(h.view)
	h.view.TitleColor = onTitleColor

	x, y := h.view.Size()
	h.table = tui.NewTableRenderer[*tui.KeybindingInfo]()
	h.table.Init(x, y, []string{"Key", "Description"}, []float64{0.2, 0.8})
	h.table.SetStyles([]color.Style{
		color.New(color.FgYellow, color.Bold),
		color.New(color.Cyan, color.OpItalic),
	})

	all := make([]*tui.KeybindingInfo, 0)
	covered := map[string]struct{}{}
	for _, k := range append(v.Keybindings, a.ui.Keybindings...) {
		if k.Key == "" {
			continue
		}
		_, exists := covered[k.Key]
		if exists {
			continue
		}

		covered[k.Key] = struct{}{}
		all = append(all, k)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Description < all[j].Description
	})

	globalTableRows := make([]*tui.TableRow[*tui.KeybindingInfo], 0)
	for _, ki := range all {
		globalTableRows = append(globalTableRows, &tui.TableRow[*tui.KeybindingInfo]{
			Cols: []string{
				ki.Key,
				ki.Description,
			},
			Value: ki,
		})
	}

	h.table.Fill(globalTableRows)

	down := func() {
		h.table.Down()
		h.show()
	}
	up := func() {
		h.table.Up()
		h.show()
	}
	prevView := a.ui.FocusedView()
	a.ui.KeyBinding(h.view).
		Set('j', "Move down", down).
		Set('k', "Move up", up).
		Set(gocui.KeyArrowDown, "Move down", down).
		Set(gocui.KeyArrowUp, "Move up", up).
		Set('g', "Go to top", func() {
			h.table.Top()
			h.show()
		}).
		Set('G', "Go to bottom", func() {
			h.table.Bottom()
			h.show()
		}).
		Set('?', "Close cheatsheet", func() {
			a.ui.DeleteView(h.view)
			if prevView != nil {
				a.ui.FocusView(prevView)
			}
		}).
		Set(gocui.KeyEsc, "Close cheatsheet", func() {
			a.ui.DeleteView(h.view)
			if prevView != nil {
				a.ui.FocusView(prevView)
			}
		})

	h.show()
	a.ui.FocusView(h.view)
}

func (h *Help) show() {
	h.view.Clear()
	h.table.Render(h.view)
}
