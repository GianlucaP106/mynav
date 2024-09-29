package ui

import (
	"fmt"
	"log"
	"mynav/pkg/core"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type subworkspaceDialog struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*core.Workspace]
	parent        *core.Workspace
}

func openSubworkspaceDialog(parent *core.Workspace) *subworkspaceDialog {
	s := &subworkspaceDialog{
		parent: parent,
	}
	screenX, screenY := tui.ScreenSize()
	s.view = tui.SetCenteredView(SubworkspaceDialog, screenX/2, screenY/2, 0)

	s.view.Title = tui.WithSurroundingSpaces("Subworkspaces")
	s.view.FrameColor = onFrameColor
	ui.styleView(s.view)

	viewX, viewY := s.view.Size()
	s.tableRenderer = tui.NewTableRenderer[*core.Workspace]()
	s.tableRenderer.InitTable(viewX, viewY, []string{
		"Name",
		"Git Remote",
		"Last Modified",
		"Tmux Session",
	}, []float64{
		0.25,
		0.25,
		0.25,
		0.25,
	})

	prevView := tui.GetFocusedView()
	s.view.KeyBinding().
		Set('j', "Move down", func() {
			s.tableRenderer.Down()
			s.render()
		}).
		Set('k', "Move up", func() {
			s.tableRenderer.Up()
			s.render()
		}).
		Set(gocui.KeyEsc, "Close dialog", func() {
			s.close()
			if prevView != nil {
				prevView.Focus()
			}
		}).
		Set('a', "Add subworkspace", func() {
			openEditorDialog(func(str string) {
				_, err := api().Workspaces.CreateSubworkspace(str, parent)
				if err != nil {
					panic(err)
				}
				s.refresh()
				s.render()
			}, func() {
			}, "Subworkspace name", smallEditorSize)
		})

	s.view.Focus()
	return nil
}

func (s *subworkspaceDialog) refresh() {
	sub, err := api().Workspaces.GetSubworkspaces(s.parent)
	if err != nil {
		log.Fatal(err)
	}

	rows := make([][]string, 0)
	for _, w := range sub {
		fmt.Println(w.Name)
		rows = append(rows, []string{
			w.Name,
			"",
			"",
			"",
		})
	}

	s.tableRenderer.FillTable(rows, sub)
}

func (s *subworkspaceDialog) close() {
	s.view.Delete()
}

func (s *subworkspaceDialog) render() {
	s.view.Clear()
	s.tableRenderer.Render(s.view)
}
