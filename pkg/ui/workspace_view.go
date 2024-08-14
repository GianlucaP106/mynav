package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/core"
	"mynav/pkg/events"
	"mynav/pkg/github"
	"mynav/pkg/system"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type WorkspacesView struct {
	view          *View
	tableRenderer *TableRenderer[*core.Workspace]
	search        string
}

var _ Viewable = new(WorkspacesView)

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
	wv.view = GetViewPosition(constants.WorkspacesViewName).Set()

	wv.view.Title = withSurroundingSpaces("Workspaces")
	StyleView(wv.View())

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
	wv.tableRenderer = NewTableRenderer[*core.Workspace]()
	wv.tableRenderer.InitTable(sizeX, sizeY, titles, proportions)

	events.AddEventListener(constants.WorkspaceChangeEventName, func(_ string) {
		wv.refresh()
		RenderView(wv)
	})

	wv.refresh()

	if selectedWorkspace := Api().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		wv.selectWorkspaceByShortPath(selectedWorkspace.ShortPath())
	}

	wv.view.KeyBinding().
		set('j', "Move down", func() {
			wv.tableRenderer.Down()
		}).
		set('k', "Move up", func() {
			wv.tableRenderer.Up()
		}).
		set(gocui.KeyEsc, "Escape search / Go back", func() {
			if wv.search != "" {
				wv.search = ""
				wv.view.Subtitle = ""
				wv.refresh()
				return
			}

			GetTopicsView().Focus()
		}).
		set('s', "See workspace information", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			OpenWorkspaceInfoDialog(curWorkspace, func() {})
		}).
		set('g', "Clone git repo", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			OpenEditorDialog(func(s string) {
				if err := Api().Core.CloneRepo(s, curWorkspace); err != nil {
					OpenToastDialogError(err.Error())
				}
			}, func() {}, "Git repo URL", Small)
		}).
		set('G', "Open browser to git repo", func() {
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
		set('/', "Search by name", func() {
			OpenEditorDialog(func(s string) {
				if s != "" {
					wv.search = s
					wv.view.Subtitle = withSurroundingSpaces("Searching: " + wv.search)
					wv.refresh()
				}
			}, func() {}, "Search", Small)
		}).
		set(gocui.KeyEnter, "Open in tmux/open in neovim", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			RunAction(func() {
				if core.IsTmuxSession() {
					Api().Core.OpenNeovimInWorkspace(curWorkspace)
				} else {
					Api().Core.CreateOrAttachTmuxSession(curWorkspace)
				}
			})
		}).
		set('v', "Open in neovim", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			RunAction(func() {
				Api().Core.OpenNeovimInWorkspace(curWorkspace)
			})
		}).
		set('t', "Open in terminal", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			RunAction(func() {
				Api().Core.OpenTerminalInWorkspace(curWorkspace)
			})
		}).
		set('m', "Move workspace to another topic", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			sd := new(*SearchListDialog[*core.Topic])
			*sd = OpenSearchListDialog(SearchDialogConfig[*core.Topic]{
				onSearch: func(s string) ([][]string, []*core.Topic) {
					rows := make([][]string, 0)
					topics := Api().Core.GetTopics().FilterByNameContaining(s)
					for _, t := range topics {
						rows = append(rows, []string{
							t.Name,
						})
					}

					return rows, topics
				},
				initial: func() ([][]string, []*core.Topic) {
					rows := make([][]string, 0)
					topics := Api().Core.GetTopics()
					for _, t := range topics {
						rows = append(rows, []string{
							t.Name,
						})
					}

					return rows, topics
				},
				onSelect: func(a *core.Topic) {
					if err := Api().Core.MoveWorkspace(curWorkspace, a); err != nil {
						OpenToastDialogError(err.Error())
						return
					}

					if *sd != nil {
						(*sd).Close()
					}

					GetTopicsView().tableRenderer.SelectRow(0)

					wv.Focus()
				},
				onSelectDescription: "Move workspace to this topic",
				searchViewTitle:     "Filter",
				tableViewTitle:      "Result",
				focusList:           true,
				tableTitles: []string{
					"Name",
				},
				tableProportions: []float64{
					1.0,
				},
			})
		}).
		set('D', "Delete a workspace", func() {
			if Api().Core.GetWorkspacesByTopicCount(GetTopicsView().getSelectedTopic()) <= 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					curWorkspace := wv.getSelectedWorkspace()
					Api().Core.DeleteWorkspace(curWorkspace)

					// HACK: same as below
					GetTopicsView().tableRenderer.SelectRow(0)
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		set('r', "Rename workspace", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			OpenEditorDialogWithDefaultValue(func(s string) {
				if err := Api().Core.RenameWorkspace(curWorkspace, s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}
			}, func() {}, "New workspace name", Small, curWorkspace.Name)
		}).
		set('e', "Add/change description", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			OpenEditorDialog(func(desc string) {
				if desc != "" {
					Api().Core.SetDescription(desc, curWorkspace)
				}
			}, func() {}, "Description", Large)
		}).
		set('a', "Create a workspace", func() {
			curTopic := GetTopicsView().getSelectedTopic()
			if curTopic == nil {
				OpenToastDialog("You must create a topic first", false, "Note", func() {})
				return
			}

			OpenEditorDialog(func(name string) {
				if _, err := Api().Core.CreateWorkspace(name, curTopic); err != nil {
					OpenToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new workspace
				// This will result in the workspace and the corresponding topic going to the top
				// because we are sorting by modifed time
				GetTopicsView().tableRenderer.SelectRow(0)
				wv.tableRenderer.SelectRow(0)
			}, func() {}, "Workspace name ", Small)
		}).
		set('X', "Kill tmux session", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if Api().Tmux.GetTmuxSessionByName(curWorkspace.Path) != nil {
				OpenConfirmationDialog(func(b bool) {
					if b {
						Api().Core.DeleteWorkspaceTmuxSession(curWorkspace)
					}
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(wv.view.keybindingInfo.toList(), func() {})
		})
}

func (wv *WorkspacesView) selectWorkspaceByShortPath(shortPath string) {
	wv.tableRenderer.SelectRowByValue(func(w *core.Workspace) bool {
		return w.ShortPath() == shortPath
	})
}

func (wv *WorkspacesView) refresh() {
	var workspaces core.Workspaces
	if selectedTopic := GetTopicsView().getSelectedTopic(); selectedTopic != nil {
		workspaces = Api().Core.GetWorkspaces().FilterByTopic(selectedTopic)
	} else {
		workspaces = make(core.Workspaces, 0)
	}

	if wv.search != "" {
		workspaces = workspaces.FilterByNameContaining(wv.search)
	}

	rows := make([][]string, 0)
	rowValues := make([]*core.Workspace, 0)
	for _, w := range workspaces {
		tmux := func() string {
			if tm := Api().Tmux.GetTmuxSessionByName(w.Path); tm != nil {
				numWindows := strconv.Itoa(tm.Windows)
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
			remote = github.TrimGithubUrl(remote)
		}

		rowValues = append(rowValues, w)
		rows = append(rows, []string{
			w.Name,
			w.GetDescription(),
			remote,
			w.GetLastModifiedTimeFormatted(),
			tmux,
		})
	}

	wv.tableRenderer.FillTable(rows, rowValues)
}

func (wv *WorkspacesView) getSelectedWorkspace() *core.Workspace {
	_, w := wv.tableRenderer.GetSelectedRow()
	if w != nil {
		return *w
	}

	return nil
}

func (wv *WorkspacesView) Render() error {
	wv.view.Clear()

	isFocused := wv.view.IsFocused()

	wv.tableRenderer.RenderWithSelectCallBack(wv.view, func(_ int, _ *TableRow[*core.Workspace]) bool {
		return isFocused
	})

	return nil
}
