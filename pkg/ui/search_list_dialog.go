package ui

import (
	"mynav/pkg/constants"

	"github.com/awesome-gocui/gocui"
)

type SearchListDialog[T any] struct {
	searchView    *View
	tableView     *View
	tableRenderer *TableRenderer[T]
}

type SearchDialogConfig[T any] struct {
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

func OpenSearchListDialog[T any](params SearchDialogConfig[T]) *SearchListDialog[T] {
	s := &SearchListDialog[T]{}

	s.searchView = SetCenteredView(constants.SearchListDialog1ViewName, 80, 3, -7)
	s.searchView.Title = params.searchViewTitle
	s.searchView.Editable = true
	s.searchView.Editor = NewSimpleEditor(func(item string) {
		rows, rowValues := params.onSearch(item)
		s.tableRenderer.FillTable(rows, rowValues)
		s.renderTable()
		s.focusList()
	}, func() {
	})

	s.tableView = SetCenteredView(constants.SearchListDialog2ViewName, 80, 10, 0)
	s.tableView.Title = params.tableViewTitle
	tableViewX, tableViewY := s.tableView.Size()

	s.tableRenderer = NewTableRenderer[T]()
	s.tableRenderer.InitTable(tableViewX, tableViewY, params.tableTitles, params.tableProportions)

	if params.initial != nil {
		rows, rowValues := params.initial()
		s.tableRenderer.FillTable(rows, rowValues)
	}

	prevView := GetFocusedView()

	s.searchView.KeyBinding().
		set(gocui.KeyEsc, "Close dialog", func() {
			s.Close()
			if prevView != nil {
				prevView.Focus()
			}
		}).
		set(gocui.KeyTab, "Toggle focus", func() {
			s.focusList()
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(s.searchView.keybindingInfo.toList(), func() {})
		})

	s.tableView.KeyBinding().
		set(gocui.KeyEsc, "Close dialog", func() {
			s.focusSearch()
		}).
		set(gocui.KeyTab, "Toggle focus", func() {
			s.focusSearch()
		}).
		set(gocui.KeyEnter, params.onSelectDescription, func() {
			_, v := s.tableRenderer.GetSelectedRow()
			params.onSelect(*v)
		}).
		set('j', "Move down", func() {
			s.tableRenderer.Down()
			s.renderTable()
		}).
		set('k', "Move up", func() {
			s.tableRenderer.Up()
			s.renderTable()
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(s.tableView.keybindingInfo.toList(), func() {})
		})

	if params.focusList {
		s.focusList()
	} else {
		s.focusSearch()
	}
	s.renderTable()

	return s
}

func (s *SearchListDialog[T]) renderTable() {
	s.tableView.Clear()
	s.tableRenderer.Render(s.tableView)
}

func (s *SearchListDialog[T]) focusSearch() {
	s.searchView.FrameColor = gocui.ColorGreen
	s.tableView.FrameColor = gocui.ColorBlue
	ToggleCursor(true)
	s.searchView.Focus()
}

func (s *SearchListDialog[T]) focusList() {
	s.searchView.FrameColor = gocui.ColorBlue
	s.tableView.FrameColor = gocui.ColorGreen
	ToggleCursor(false)
	s.tableView.Focus()
}

func (s *SearchListDialog[T]) Close() {
	ToggleCursor(false)
	s.searchView.Delete()
	s.tableView.Delete()
}
