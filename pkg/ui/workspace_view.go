package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tmux"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type WorkspacesView struct {
	view          *View
	tableRenderer *TableRenderer
	search        string
	workspaces    core.Workspaces
}

const WorkspacesViewName = "WorkspacesView"

var _ Viewable = new(WorkspacesView)

func NewWorkspcacesView() *WorkspacesView {
	return &WorkspacesView{}
}

func GetWorkspacesView() *WorkspacesView {
	return GetViewable[*WorkspacesView]()
}

func FocusWorkspacesView() {
	FocusView(WorkspacesViewName)
}

func (wv *WorkspacesView) View() *View {
	return wv.view
}

func (wv *WorkspacesView) Init() {
	wv.view = SetViewLayout(WorkspacesViewName)

	wv.view.Title = withSurroundingSpaces("Workspaces")
	wv.view.TitleColor = gocui.ColorBlue
	wv.view.FrameColor = gocui.ColorBlue

	sizeX, sizeY := wv.view.Size()

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
		FocusTmuxView()
	}

	moveLeft := func() {
		FocusTopicsView()
	}

	KeyBinding(wv.view.Name()).
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
				wv.view.Subtitle = ""
				wv.refreshWorkspaces()
				return
			}

			FocusTopicsView()
		}).
		set('s', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			OpenWorkspaceInfoDialog(curWorkspace, func() {})
		}).
		set('g', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			OpenEditorDialog(func(s string) {
				if err := Api().Core.CloneRepo(s, curWorkspace); err != nil {
					OpenToastDialogError(err.Error())
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
				OpenToastDialogError(err.Error())
			}
		}).
		set('/', func() {
			OpenEditorDialog(func(s string) {
				if s != "" {
					wv.search = s
					wv.view.Subtitle = withSurroundingSpaces("Searching: " + wv.search)
					wv.refreshWorkspaces()
				}
			}, func() {}, "Search", Small)
		}).
		setWithQuit(gocui.KeyEnter, func() bool {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return false
			}

			if tmux.IsTmuxSession() {
				SetAction(Api().Core.GetWorkspaceNvimCmd(curWorkspace))
				return true
			}

			command := Api().Core.GetCreateOrAttachTmuxSessionCmd(curWorkspace)
			SetAction(command)

			return true
		}).
		setWithQuit('v', func() bool {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return false
			}

			Api().Core.SetSelectedWorkspace(curWorkspace)
			SetAction(system.GetNvimCmd(curWorkspace.Path))
			return true
		}).
		setWithQuit('m', func() bool {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return false
			}

			openTermCmd, err := system.GetOpenTerminalCmd(curWorkspace.Path)
			if err != nil {
				OpenToastDialogError(err.Error())
				return false
			}

			SetAction(openTermCmd)
			return true
		}).
		set('D', func() {
			if Api().Core.GetWorkspacesByTopicCount(GetTopicsView().getSelectedTopic()) <= 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					curWorkspace := wv.getSelectedWorkspace()
					Api().Core.DeleteWorkspace(curWorkspace)

					// HACK: same as below
					GetTopicsView().tableRenderer.SetSelectedRow(0)
					RefreshAllData()
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		set('r', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			OpenEditorDialogWithDefaultValue(func(s string) {
				if err := Api().Core.RenameWorkspace(curWorkspace, s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}
				wv.syncWorkspacesToTable()
			}, func() {}, "New workspace name", Small, curWorkspace.Name)
		}).
		set('e', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			OpenEditorDialog(func(desc string) {
				if desc != "" {
					Api().Core.SetDescription(desc, curWorkspace)
					wv.syncWorkspacesToTable()
				}
			}, func() {}, "Description", Large)
		}).
		set('a', func() {
			curTopic := GetTopicsView().getSelectedTopic()
			OpenEditorDialog(func(name string) {
				if _, err := Api().Core.CreateWorkspace(name, curTopic); err != nil {
					OpenToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new workspace
				// This will result in the workspace and the corresponding topic going to the top
				// because we are sorting by modifed time
				GetTmuxSessionView().tableRenderer.SetSelectedRow(0)
				wv.tableRenderer.SetSelectedRow(0)
				RefreshAllData()
			}, func() {}, "Workspace name ", Small)
		}).
		set('X', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if Api().Tmux.GetTmuxSessionByName(curWorkspace.Path) != nil {
				OpenConfirmationDialog(func(b bool) {
					if b {
						Api().Core.DeleteWorkspaceTmuxSession(curWorkspace)
						RefreshAllData()
					}
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		set('?', func() {
			OpenHelpView(workspaceKeyBindings, func() {})
		})
}

func (wv *WorkspacesView) selectWorkspaceByShortPath(shortPath string) {
	for idx, w := range wv.workspaces {
		if w.ShortPath() == shortPath {
			wv.tableRenderer.SetSelectedRow(idx)
		}
	}
}

func (wv *WorkspacesView) refreshWorkspaces() {
	var workspaces core.Workspaces
	if selectedTopic := GetTopicsView().getSelectedTopic(); selectedTopic != nil {
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

func (wv *WorkspacesView) Render() error {
	wv.view.Clear()

	// TODO: refact
	isFocused := false
	if v := GetFocusedView(); v != nil && v.Name() == wv.view.Name() {
		isFocused = true
	}

	wv.tableRenderer.RenderWithSelectCallBack(wv.view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	return nil
}
