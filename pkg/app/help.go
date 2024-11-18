package app

import (
	"fmt"
	"mynav/pkg/tui"
	"sort"

	"github.com/gookit/color"
)

type Help struct {
	view *tui.View

	// tables for table and local keys
	table *tui.TableRenderer[*tui.KeybindingInfo]
	// local  *tui.TableRenderer[*tui.KeybindingInfo]
}

func help(v *tui.View) {
	h := &Help{}
	h.view = a.ui.SetCenteredView(HelpDialog, 80, 20, 0)
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
		return all[i].Key < all[j].Key
	})

	globalRows := make([][]string, 0)
	for _, ki := range all {
		globalRows = append(globalRows, []string{
			ki.Key,
			ki.Description,
		})
	}
	h.table.Fill(globalRows, all)

	a.ui.KeyBinding(h.view).
		Set('?', "Close cheatsheet", func() {
			a.ui.DeleteView(h.view)
		}).
		Set('j', "Move down", func() {
			h.table.Down()
			h.show()
		}).
		Set('k', "Move up", func() {
			h.table.Up()
			h.show()
		})

	h.show()
	a.ui.FocusView(h.view)
}

func (h *Help) show() {
	h.view.Clear()
	h.table.Render(h.view)
}
