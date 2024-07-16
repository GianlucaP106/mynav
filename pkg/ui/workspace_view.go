package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/git"
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

var _ Viewable = new(WorkspacesView)

const WorkspacesViewName = "WorkspacesView"

func NewWorkspcacesView() *WorkspacesView {
	return &WorkspacesView{}
}

func GetWorkspacesView() *WorkspacesView {
	return GetViewable[*WorkspacesView]()
}

func (wv *WorkspacesView) View() *View {
	return wv.view
}

func (wv *WorkspacesView) Focus() {
	FocusView(wv.View().Name())
}

func (wv *WorkspacesView) Init() {
	wv.view = GetViewPosition(WorkspacesViewName).Set()

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

	moveLeft := func() {
		GetTopicsView().Focus()
	}

	wv.view.KeyBinding().
		set('j', func() {
			wv.tableRenderer.Down()
		}, "Move down").
		set('k', func() {
			wv.tableRenderer.Up()
		}, "Move up").
		set(gocui.KeyArrowLeft, moveLeft, "Focus topic view").
		set(gocui.KeyCtrlH, moveLeft, "Focus topic view").
		set(gocui.KeyEsc, func() {
			if wv.search != "" {
				wv.search = ""
				wv.view.Subtitle = ""
				wv.refreshWorkspaces()
				return
			}

			GetTopicsView().Focus()
		}, "Escape search / Go back").
		set('s', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			OpenWorkspaceInfoDialog(curWorkspace, func() {})
		}, "See workspace information").
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
		}, "Clone git repo").
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
		}, "Open browser to git repo").
		set('/', func() {
			OpenEditorDialog(func(s string) {
				if s != "" {
					wv.search = s
					wv.view.Subtitle = withSurroundingSpaces("Searching: " + wv.search)
					wv.refreshWorkspaces()
				}
			}, func() {}, "Search", Small)
		}, "Search by name").
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
		}, "Open in tmux/open in neovim").
		setWithQuit('v', func() bool {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return false
			}

			Api().Core.SetSelectedWorkspace(curWorkspace)
			SetAction(system.GetNvimCmd(curWorkspace.Path))
			return true
		}, "Open in neovim").
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
		}, "Open in terminal").
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
		}, "Delete a workspace").
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
		}, "Rename workspace").
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
		}, "Add/change description").
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
				GetTopicsView().tableRenderer.SetSelectedRow(0)
				wv.tableRenderer.SetSelectedRow(0)
				RefreshAllData()
			}, func() {}, "Workspace name ", Small)
		}, "Create a workspace").
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
		}, "Kill tmux session").
		set('?', func() {
			OpenHelpView(wv.view.keybindingInfo.toList(), func() {})
		}, "Toggle cheatsheet")
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

		remote, err := w.GetGitRemote()
		if err != nil {
			OpenToastDialogError(err.Error())
			return
		}

		if remote != "" {
			remote = git.TrimGithubUrl(remote)
		}

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

	isFocused := wv.view.IsFocused()

	wv.tableRenderer.RenderWithSelectCallBack(wv.view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	return nil
}
