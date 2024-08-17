package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type searchListDialog[T any] struct {
	searchView    *tui.View
	tableView     *tui.View
	tableRenderer *tui.TableRenderer[T]
}

type searchDialogConfig[T any] struct {
	onSearch            func(s string) ([][]string, []T)
	onSelect            func(a T)
	initial             func() ([][]string, []T)
	onSelectDescription string
	searchViewTitle     string
	tableViewTitle      string
	tableTitles         []string
	tableProportions    []float64
	focusList           bool
}

func openSearchListDialog[T any](params searchDialogConfig[T]) *searchListDialog[T] {
	s := &searchListDialog[T]{}

	s.searchView = tui.SetCenteredView(constants.SearchListDialog1ViewName, 80, 3, -7)
	s.searchView.Title = params.searchViewTitle
	s.searchView.Subtitle = tui.WithSurroundingSpaces("<Tab> to toggle focus")
	s.searchView.Editable = true
	s.searchView.Editor = tui.NewSimpleEditor(func(item string) {
		rows, rowValues := params.onSearch(item)
		s.tableRenderer.FillTable(rows, rowValues)
		s.renderTable()
		s.focusList()
	}, func() {
	})

	tui.StyleView(s.searchView)

	s.tableView = tui.SetCenteredView(constants.SearchListDialog2ViewName, 80, 10, 0)
	s.tableView.Title = params.tableViewTitle
	tableViewX, tableViewY := s.tableView.Size()
	tui.StyleView(s.tableView)

	s.tableRenderer = tui.NewTableRenderer[T]()
	s.tableRenderer.InitTable(tableViewX, tableViewY, params.tableTitles, params.tableProportions)

	if params.initial != nil {
		rows, rowValues := params.initial()
		s.tableRenderer.FillTable(rows, rowValues)
	}

	prevView := tui.GetFocusedView()

	s.searchView.KeyBinding().
		Set(gocui.KeyEsc, "Close dialog", func() {
			s.close()
			if prevView != nil {
				prevView.Focus()
			}
		}).
		Set(gocui.KeyTab, "Toggle focus", func() {
			s.focusList()
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(s.searchView.GetKeybindings(), func() {})
		})

	s.tableView.KeyBinding().
		Set(gocui.KeyEsc, "Close dialog", func() {
			s.focusSearch()
		}).
		Set(gocui.KeyTab, "Toggle focus", func() {
			s.focusSearch()
		}).
		Set(gocui.KeyEnter, params.onSelectDescription, func() {
			_, v := s.tableRenderer.GetSelectedRow()
			params.onSelect(*v)
		}).
		Set('j', "Move down", func() {
			s.tableRenderer.Down()
			s.renderTable()
		}).
		Set('k', "Move up", func() {
			s.tableRenderer.Up()
			s.renderTable()
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(s.tableView.GetKeybindings(), func() {})
		})

	if params.focusList {
		s.focusList()
	} else {
		s.focusSearch()
	}
	s.renderTable()

	return s
}

func (s *searchListDialog[T]) renderTable() {
	s.tableView.Clear()
	s.tableRenderer.Render(s.tableView)
}

func (s *searchListDialog[T]) focusSearch() {
	s.searchView.FrameColor = tui.OnFrameColor
	s.tableView.FrameColor = tui.OffFrameColor
	tui.ToggleCursor(true)
	s.searchView.Focus()
}

func (s *searchListDialog[T]) focusList() {
	s.searchView.FrameColor = tui.OffFrameColor
	s.tableView.FrameColor = tui.OnFrameColor
	tui.ToggleCursor(false)
	s.tableView.Focus()
}

func (s *searchListDialog[T]) close() {
	tui.ToggleCursor(false)
	s.searchView.Delete()
	s.tableView.Delete()
}
