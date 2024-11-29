package app

import (
	"fmt"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

// Search is dialog that provides a list and editor for searching purposes.
type Search[T any] struct {
	searchView *tui.View
	tableView  *tui.View
	table      *tui.TableRenderer[T]
}

// Params of the search dialog.
type searchDialogConfig[T any] struct {
	onSearch            func(s string) ([][]string, []T)
	onSelect            func(a T)
	initial             func() ([][]string, []T)
	onSelectDescription string
	searchViewTitle     string
	tableViewTitle      string
	tableTitles         []string
	colStyles           []color.Style
	tableProportions    []float64
	focusList           bool
}

// Opens the search dialog with the given params.
func search[T any](params searchDialogConfig[T]) *Search[T] {
	// build search view
	s := &Search[T]{}
	s.searchView = a.ui.SetCenteredView(SearchListDialog1View, 80, 3, -7)
	s.searchView.Title = fmt.Sprintf(" %s ", params.searchViewTitle)
	s.searchView.Subtitle = " <Enter> to filter "
	s.searchView.Editable = true
	a.styleView(s.searchView)
	s.searchView.Editor = tui.NewSimpleEditor(func(item string) {
		rows, rowValues := params.onSearch(item)
		s.table.Fill(rows, rowValues)
		s.renderTable()
		s.focusList()
	}, func() {
	})

	// build table view
	s.tableView = a.ui.SetCenteredView(SearchListDialog2View, 80, 10, 0)
	s.tableView.Title = fmt.Sprintf(" %s ", params.tableViewTitle)
	s.tableView.Subtitle = " <Tab> to toggle focus "
	tableViewX, tableViewY := s.tableView.Size()
	a.styleView(s.tableView)

	// build and fill table
	s.table = tui.NewTableRenderer[T]()
	s.table.Init(tableViewX, tableViewY, params.tableTitles, params.tableProportions)
	if params.colStyles != nil {
		s.table.SetStyles(params.colStyles)
	}
	if params.initial != nil {
		rows, rowValues := params.initial()
		s.table.Fill(rows, rowValues)
	}

	// keybindings
	prevView := a.ui.FocusedView()
	a.ui.KeyBinding(s.searchView).
		Set(gocui.KeyEsc, "Close dialog", func() {
			s.close()
			if prevView != nil {
				a.ui.FocusView(prevView)
			}
		}).
		Set(gocui.KeyTab, "Toggle focus", func() {
			s.focusList()
		})

	a.ui.KeyBinding(s.tableView).
		Set(gocui.KeyEsc, "Close dialog", func() {
			s.close()
			if prevView != nil {
				a.ui.FocusView(prevView)
			}
		}).
		Set(gocui.KeyTab, "Toggle focus", func() {
			s.focusSearch()
		}).
		Set(gocui.KeyEnter, params.onSelectDescription, func() {
			_, v := s.table.SelectedRow()
			if v != nil {
				params.onSelect(*v)
			}
		}).
		Set('j', "Move down", func() {
			s.table.Down()
			s.renderTable()
		}).
		Set(gocui.KeyArrowDown, "Move down", func() {
			s.table.Down()
			s.renderTable()
		}).
		Set('k', "Move up", func() {
			s.table.Up()
			s.renderTable()
		}).
		Set(gocui.KeyArrowUp, "Move up", func() {
			s.table.Up()
			s.renderTable()
		}).
		Set('?', "Toggle cheatsheet", func() {
			help(s.tableView)
		})

	if params.focusList {
		s.focusList()
	} else {
		s.focusSearch()
	}
	s.renderTable()

	return s
}

func (s *Search[T]) renderTable() {
	s.tableView.Clear()
	s.table.Render(s.tableView)
}

func (s *Search[T]) focusSearch() {
	s.searchView.FrameColor = onFrameColor
	s.tableView.FrameColor = offFrameColor
	a.ui.Cursor = true
	a.ui.FocusView(s.searchView)
}

func (s *Search[T]) focusList() {
	s.searchView.FrameColor = offFrameColor
	s.tableView.FrameColor = onFrameColor
	a.ui.Cursor = false
	a.ui.FocusView(s.tableView)
}

func (s *Search[T]) close() {
	a.ui.Cursor = false
	a.ui.DeleteView(s.searchView)
	a.ui.DeleteView(s.tableView)
}
