package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/GianlucaP106/mynav/pkg/core"
	"github.com/GianlucaP106/mynav/pkg/tui"
	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

// Search is dialog that provides a list and editor for searching purposes.
type Search[T any] struct {
	previewEnabled bool
	bgView         *tui.View
	searchView     *tui.View
	tableView      *tui.View
	previewView    *Preview
	table          *tui.TableRenderer[T]
}

// Params of the search dialog.
type SearchDialogConfig[T any] struct {
	onSearch            func(s string) []*tui.TableRow[T]
	onSelect            func(a T)
	onType              func(string) []*tui.TableRow[T]
	initial             func() []*tui.TableRow[T]
	onTypePreview       func(*tui.TableRow[T]) *core.Session
	onSelectDescription string
	searchViewTitle     string
	tableViewTitle      string
	tableTitles         []string
	colStyles           []color.Style
	tableProportions    []float64
	focusList           bool
	enablePreview       bool
}

// Opens the search dialog with the given params.
func search[T any](params SearchDialogConfig[T]) *Search[T] {
	s := &Search[T]{}
	screenX, _ := a.ui.Size()

	s.previewEnabled = params.enablePreview && screenX > 160

	sizeX := 100
	horizontalOffset := 0
	if s.previewEnabled {
		horizontalOffset = -50
		sizeX = 50
	}

	bgSizeX := 100
	if s.previewEnabled {
		bgSizeX = 154
	}
	s.bgView = a.ui.SetCenteredView(SearchListDialogBgView, bgSizeX, 36, 0, 0)
	s.bgView.Frame = false

	s.searchView = a.ui.SetCenteredView(SearchListDialog1View, sizeX, 3, -16, horizontalOffset)
	s.searchView.Title = fmt.Sprintf(" %s ", params.searchViewTitle)
	s.searchView.Subtitle = " <Tab> to toggle focus "
	s.searchView.Editable = true
	a.styleView(s.searchView)

	var onType func(s string) = nil
	if params.onType != nil {
		onType = func(search string) {
			a.worker.DebounceLoad(func() {
				rows := params.onType(search)
				if s.previewEnabled && params.onTypePreview != nil {
					if len(rows) > 0 {
						session := params.onTypePreview(rows[0])
						s.previewView.setSession(session)
					}
				}
				s.table.Fill(rows)
			}, func() {
				a.ui.Update(func() {
					s.renderTable()
					if s.previewEnabled {
						s.previewView.render()
					}
				})
			}, func() {})
		}
	}
	s.searchView.Editor = tui.NewSimpleEditor(func(item string) {
		s.table.Fill(params.onSearch(item))
		s.renderTable()
		s.focusList()
	}, func() {
	}, onType)

	// build table view
	s.tableView = a.ui.SetCenteredView(SearchListDialog2View, sizeX, 31, 2, horizontalOffset)
	s.tableView.Title = fmt.Sprintf(" %s ", params.tableViewTitle)

	tableViewX, tableViewY := s.tableView.Size()
	a.styleView(s.tableView)

	// build and fill table
	s.table = tui.NewTableRenderer[T]()
	s.table.Init(tableViewX, tableViewY, params.tableTitles, params.tableProportions)
	if params.colStyles != nil {
		s.table.SetStyles(params.colStyles)
	}

	if params.initial != nil {
		s.table.Fill(params.initial())
	}

	if s.previewEnabled {
		s.previewView = newPreview()
		s.previewView.init(a.ui.SetCenteredView(SearchListDialog3View, 100, 34, 0, 26))
	}

	update := func() {
		a.worker.Queue(func() {
			if params.onTypePreview != nil {
				_, row := s.table.SelectedRow()
				session := params.onTypePreview(row)
				s.previewView.setSession(session)
			}
			a.ui.Update(func() {
				s.previewView.render()
			})
		})
		s.renderTable()
	}

	up := func() {
		s.table.Up()
		update()
	}

	down := func() {
		s.table.Down()
		update()
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
		Set(gocui.KeyCtrlJ, "Move down list", down).
		Set(gocui.KeyCtrlK, "Move up list", up).
		Set(gocui.KeyCtrlN, "Move down list", down).
		Set(gocui.KeyCtrlP, "Move up list", up).
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
				params.onSelect(v.Value)
			}
		}).
		Set('j', "Move down", down).
		Set(gocui.KeyArrowDown, "Move down", down).
		Set('k', "Move up", up).
		Set(gocui.KeyArrowUp, "Move up", up).
		Set('g', "Go to top", func() {
			s.table.Top()
			update()
		}).
		Set('G', "Go to bottom", func() {
			s.table.Bottom()
			update()
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
	row, _ := s.table.SelectedRow()
	size := s.table.Size()
	s.tableView.Subtitle = fmt.Sprintf(" %d / %d ", min(row+1, size), size)
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
	a.ui.DeleteView(s.bgView)
	if s.previewEnabled {
		s.previewView.teardown()
	}
}

type SearchItem struct {
	workspace *core.Workspace
	session   *core.Session
}

type GlobalSearch struct{}

func newGlobalSearch() *GlobalSearch {
	return &GlobalSearch{}
}

func (g *GlobalSearch) init() {
	useFzf := core.IsFzfInstalled()
	if !useFzf {
		toast("install fzf it for a better experience", toastWarn)
	}

	const workspacePrefix = "--workspace--"
	const sessionPrefix = "--session--"

	allNames := []string{}
	allWorkspaces := a.api.AllWorkspaces().Sorted()
	for _, w := range allWorkspaces {
		allNames = append(allNames, fmt.Sprintf("%s%s", workspacePrefix, w.ShortPath()))
	}

	allSessions := a.api.AllSessions()
	for _, s := range allSessions {
		allNames = append(allNames, fmt.Sprintf("%s%s", sessionPrefix, s.Name))
	}

	searchFor := func(s string) []*tui.TableRow[SearchItem] {
		foundItems := make([]SearchItem, 0)
		if useFzf {
			found := core.FuzzyFind(allNames, s)
			for _, item := range found {
				if item == "" {
					continue
				}

				switch {
				case strings.HasPrefix(item, workspacePrefix):
					i := strings.TrimPrefix(item, workspacePrefix)
					w := a.api.Workspace(i)
					if w != nil {
						foundItems = append(foundItems, SearchItem{
							workspace: w,
						})
					}
				case strings.HasPrefix(item, sessionPrefix):
					i := strings.TrimPrefix(item, sessionPrefix)
					session, _ := a.api.SessionByName(i)
					if session != nil {
						foundItems = append(foundItems, SearchItem{
							session: session,
						})
					}
				}
			}
		} else {
			for _, workspace := range allWorkspaces {
				if strings.Contains(workspace.Name, s) {
					foundItems = append(foundItems, SearchItem{
						workspace: workspace,
					})
				}
			}
			for _, session := range allSessions {
				if strings.Contains(session.Name, s) {
					foundItems = append(foundItems, SearchItem{
						session: session,
					})
				}
			}
		}

		tableRows := make([]*tui.TableRow[SearchItem], 0)
		for _, item := range foundItems {
			cols := []string{}
			styles := []color.Style{}
			switch {
			case item.session != nil:
				name := item.session.Name
				if parts := strings.Split(name, "/"); len(parts) > 3 {
					lastParts := parts[len(parts)-3:]
					lastParts = append([]string{".../"}, lastParts...)
					name = filepath.Join(lastParts...)
				}

				cols = []string{
					name,
					"Session",
				}
				styles = []color.Style{
					workspaceNameColor,
					sessionMarkerColor,
				}
			case item.workspace != nil:
				cols = []string{
					item.workspace.ShortPath(),
					"Workspace",
				}
				styles = []color.Style{
					workspaceNameColor,
					alternateSessionMarkerColor,
				}
			}
			tableRows = append(tableRows, &tui.TableRow[SearchItem]{
				Cols:   cols,
				Value:  item,
				Styles: styles,
			})
		}
		return tableRows
	}

	sd := new(*Search[SearchItem])
	*sd = search(SearchDialogConfig[SearchItem]{
		onType: searchFor,
		onTypePreview: func(tr *tui.TableRow[SearchItem]) *core.Session {
			if tr == nil {
				return nil
			}

			if tr.Value.session != nil {
				return tr.Value.session
			}

			return a.api.Session(tr.Value.workspace)
		},
		onSearch: searchFor,
		onSelect: func(item SearchItem) {
			if *sd != nil {
				(*sd).close()
			}
			switch {
			case item.session != nil:
				a.sessions.selectSession(item.session)
				a.sessions.focus()
			case item.workspace != nil:
				a.topics.selectTopic(item.workspace.Topic)
				a.workspaces.refresh()
				a.workspaces.selectWorkspace(item.workspace)
				a.workspaces.focus()
			}
		},
		onSelectDescription: "Go to item",
		searchViewTitle:     "Search",
		tableViewTitle:      "Result",
		tableTitles: []string{
			"Name",
			"Type",
		}, tableProportions: []float64{
			0.6,
			0.4,
		},
		colStyles: []color.Style{
			workspaceNameColor,
			sessionMarkerColor,
		},
		enablePreview: true,
	})
}
