package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/system"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

const WorkspacesViewName = "WorkspacesView"

type WorkspacesView struct {
	tableRenderer *TableRenderer
	tv            *TopicsView
	search        string
	workspaces    core.Workspaces
}

var _ View = &WorkspacesView{}

func newWorkspacesView(tv *TopicsView) *WorkspacesView {
	workspace := &WorkspacesView{
		tv: tv,
	}
	return workspace
}

func (wv *WorkspacesView) selectWorkspaceByShortPath(shortPath string) {
	for idx, w := range wv.workspaces {
		if w.ShortPath() == shortPath {
			wv.tv.tableRenderer.SetSelectedRow(idx)
		}
	}
}

func (wv *WorkspacesView) refreshWorkspaces() {
	tv := wv.tv
	var workspaces core.Workspaces
	if selectedTopic := tv.getSelectedTopic(); selectedTopic != nil {
		workspaces = Api().Core.GetWorkspaces().ByTopic(selectedTopic)
	} else {
		workspaces = make(core.Workspaces, 0)
	}

	if wv.search != "" {
		workspaces = workspaces.FilterByNameContaining(wv.search)
	}

	wv.workspaces = workspaces
	wv.syncWorkspacesToTable()
}

func (wv *WorkspacesView) syncWorkspacesToTable() {
	rows := make([][]string, 0)
	for _, w := range wv.workspaces {
		tmux := func() string {
			if tm := Api().Tmux.GetTmuxSessionByName(w.Path); tm != nil {
				numWindows := strconv.Itoa(tm.NumWindows)
				// TODO:add color to tmux
				return numWindows + " - tmux"
			}

			return ""
		}()

		// TODO: handle
		remote, _ := w.GetGitRemote()

		rows = append(rows, []string{
			w.Name,
			w.GetDescription(),
			remote,
			w.GetLastModifiedTimeFormatted(),
			tmux,
		})
	}

	wv.tableRenderer.FillTable(rows)
}

func (wv *WorkspacesView) getSelectedWorkspace() *core.Workspace {
	idx := wv.tableRenderer.GetSelectedRowIndex()
	if idx >= len(wv.workspaces) || idx < 0 {
		return nil
	}
	return wv.workspaces[idx]
}

func (wv *WorkspacesView) RequiresManager() bool {
	return false
}

func (wv *WorkspacesView) Name() string {
	return WorkspacesViewName
}

func (wv *WorkspacesView) Init(ui *UI) {
	if GetInternalView(wv.Name()) != nil {
		return
	}

	view := SetViewLayout(wv.Name())

	view.Title = withSurroundingSpaces("Workspaces")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorBlue

	// FIX: this is not the best place to but this
	if Api().Core.GetSelectedWorkspace() != nil {
		ui.FocusWorkspacesView()
	} else {
		ui.FocusTopicsView()
	}

	sizeX, sizeY := view.Size()

	titles := []string{
		"Name",
		"Description",
		"Git remote",
		"Last Modified",
		"Tmux Session",
	}
	proportions := []float64{
		0.2,
		0.2,
		0.2,
		0.2,
		0.2,
	}
	wv.tableRenderer = NewTableRenderer()
	wv.tableRenderer.InitTable(sizeX, sizeY, titles, proportions)
	wv.refreshWorkspaces()

	if selectedWorkspace := Api().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		wv.selectWorkspaceByShortPath(selectedWorkspace.ShortPath())
	}

	moveDown := func() {
		ui.FocusTmuxView()
	}

	moveLeft := func() {
		ui.FocusTopicsView()
	}

	KeyBinding(wv.Name()).
		set('j', func() {
			wv.tableRenderer.Down()
		}).
		set('k', func() {
			wv.tableRenderer.Up()
		}).
		set(gocui.KeyArrowDown, moveDown).
		set(gocui.KeyCtrlJ, moveDown).
		set(gocui.KeyArrowLeft, moveLeft).
		set(gocui.KeyCtrlH, moveLeft).
		set(gocui.KeyEsc, func() {
			if wv.search != "" {
				wv.search = ""
				wv.refreshWorkspaces()
				return
			}

			ui.FocusTopicsView()
		}).
		set('s', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			GetDialog[*WorkspaceInfoDialog](ui).Open(curWorkspace, func() {})
		}).
		set('g', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().Core.CloneRepo(s, curWorkspace); err != nil {
					GetDialog[*ToastDialog](ui).OpenError(err.Error())
				}
				wv.syncWorkspacesToTable()
			}, func() {}, "Git repo URL", Small)
		}).
		set('G', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if curWorkspace.GitRemote == nil {
				return
			}

			if err := system.OpenBrowser(*curWorkspace.GitRemote); err != nil {
				GetDialog[*ToastDialog](ui).OpenError(err.Error())
			}
		}).
		set('/', func() {
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				wv.search = s
				wv.refreshWorkspaces()
			}, func() {}, "Search", Small)
		}).
		setKeybinding(wv.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			if system.IsTmuxSession() {
				ui.setAction(Api().Core.GetWorkspaceNvimCmd(curWorkspace))
				return gocui.ErrQuit
			}

			command := Api().Core.GetCreateOrAttachTmuxSessionCmd(curWorkspace)
			ui.setAction(command)

			return gocui.ErrQuit
		}).
		setKeybinding(wv.Name(), 'v', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			Api().Core.SetSelectedWorkspace(curWorkspace)
			ui.setAction(system.GetNvimCmd(curWorkspace.Path))
			return gocui.ErrQuit
		}).
		setKeybinding(wv.Name(), 'm', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			openTermCmd, err := system.GetOpenTerminalCmd(curWorkspace.Path)
			if err != nil {
				GetDialog[*ToastDialog](ui).OpenError(err.Error())
				return nil
			}

			ui.setAction(openTermCmd)
			return gocui.ErrQuit
		}).
		set('D', func() {
			if Api().Core.GetWorkspacesByTopicCount(wv.tv.getSelectedTopic()) <= 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					curWorkspace := wv.getSelectedWorkspace()
					Api().Core.DeleteWorkspace(curWorkspace)

					// HACK: same as below
					wv.tv.tableRenderer.SetSelectedRow(0)
					ui.RefreshMainView()
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		set('r', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().Core.RenameWorkspace(curWorkspace, s); err != nil {
					GetDialog[*ToastDialog](ui).OpenError(err.Error())
					return
				}
				wv.syncWorkspacesToTable()
			}, func() {}, "New workspace name", Small)
		}).
		set('e', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(desc string) {
				if desc != "" {
					Api().Core.SetDescription(desc, curWorkspace)
					wv.syncWorkspacesToTable()
				}
			}, func() {}, "Description", Large)
		}).
		set('a', func() {
			tv := wv.tv
			curTopic := tv.getSelectedTopic()
			GetDialog[*EditorDialog](ui).Open(func(name string) {
				if _, err := Api().Core.CreateWorkspace(name, curTopic); err != nil {
					GetDialog[*ToastDialog](ui).OpenError(err.Error())
					return
				}

				// HACK: when there a is a new workspace
				// This will result in the workspace and the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.tableRenderer.SetSelectedRow(0)
				wv.tableRenderer.SetSelectedRow(0)
				ui.RefreshMainView()
			}, func() {}, "Workspace name ", Small)
		}).
		set('X', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if Api().Tmux.GetTmuxSessionByName(curWorkspace.Path) != nil {
				GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
					if b {
						Api().Core.DeleteWorkspaceTmuxSession(curWorkspace)
						ui.RefreshMainView()
					}
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(workspaceKeyBindings, func() {})
		})
}

func (wv *WorkspacesView) Render(ui *UI) error {
	view := GetInternalView(wv.Name())
	if view == nil {
		wv.Init(ui)
		view = GetInternalView(wv.Name())
	}

	if wv.search != "" {
		view.Subtitle = withSurroundingSpaces("Searching: " + wv.search)
	} else {
		view.Subtitle = ""
	}

	view.Clear()
	isFocused := GetFocusedView().Name() == wv.Name()
	wv.tableRenderer.RenderWithSelectCallBack(view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	return nil
}
