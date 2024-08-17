package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/core"
	"mynav/pkg/events"
	"mynav/pkg/github"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type workspacesView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*core.Workspace]
	search        string
}

var _ viewable = new(workspacesView)

func newWorkspcacesView() *workspacesView {
	return &workspacesView{}
}

func getWorkspacesView() *workspacesView {
	return getViewable[*workspacesView]()
}

func (wv *workspacesView) getView() *tui.View {
	return wv.view
}

func (wv *workspacesView) Focus() {
	focusView(wv.getView().Name())
}

func (wv *workspacesView) init() {
	wv.view = GetViewPosition(constants.WorkspacesViewName).Set()

	wv.view.Title = tui.WithSurroundingSpaces("Workspaces")
	tui.StyleView(wv.getView())

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
	wv.tableRenderer = tui.NewTableRenderer[*core.Workspace]()
	wv.tableRenderer.InitTable(sizeX, sizeY, titles, proportions)

	events.AddEventListener(constants.WorkspaceChangeEventName, func(_ string) {
		wv.refresh()
		renderView(wv)
	})

	wv.refresh()

	if selectedWorkspace := getApi().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		wv.selectWorkspaceByShortPath(selectedWorkspace.ShortPath())
	}

	wv.view.KeyBinding().
		Set('j', "Move down", func() {
			wv.tableRenderer.Down()
		}).
		Set('k', "Move up", func() {
			wv.tableRenderer.Up()
		}).
		Set(gocui.KeyEsc, "Escape search / Go back", func() {
			if wv.search != "" {
				wv.search = ""
				wv.view.Subtitle = ""
				wv.refresh()
				return
			}

			getTopicsView().Focus()
		}).
		Set('s', "See workspace information", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			openWorkspaceInfoDialog(curWorkspace, func() {})
		}).
		Set('g', "Clone git repo", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			openEditorDialog(func(s string) {
				if err := getApi().Core.CloneRepo(s, curWorkspace); err != nil {
					openToastDialogError(err.Error())
				}
			}, func() {}, "Git repo URL", smallEditorSize)
		}).
		Set('G', "Open browser to git repo", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if curWorkspace.GitRemote == nil {
				return
			}

			if err := system.OpenBrowser(*curWorkspace.GitRemote); err != nil {
				openToastDialogError(err.Error())
			}
		}).
		Set('/', "Search by name", func() {
			openEditorDialog(func(s string) {
				if s != "" {
					wv.search = s
					wv.view.Subtitle = tui.WithSurroundingSpaces("Searching: " + wv.search)
					wv.refresh()
				}
			}, func() {}, "Search", smallEditorSize)
		}).
		Set(gocui.KeyEnter, "Open in tmux/open in neovim", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			tui.RunAction(func() {
				if core.IsTmuxSession() {
					getApi().Core.OpenNeovimInWorkspace(curWorkspace)
				} else {
					getApi().Core.CreateOrAttachTmuxSession(curWorkspace)
				}
			})
		}).
		Set('v', "Open in neovim", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			tui.RunAction(func() {
				getApi().Core.OpenNeovimInWorkspace(curWorkspace)
			})
		}).
		Set('t', "Open in terminal", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			tui.RunAction(func() {
				getApi().Core.OpenTerminalInWorkspace(curWorkspace)
			})
		}).
		Set('m', "Move workspace to another topic", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			sd := new(*searchListDialog[*core.Topic])
			*sd = openSearchListDialog(searchDialogConfig[*core.Topic]{
				onSearch: func(s string) ([][]string, []*core.Topic) {
					rows := make([][]string, 0)
					topics := getApi().Core.GetTopics().FilterByNameContaining(s)
					for _, t := range topics {
						rows = append(rows, []string{
							t.Name,
						})
					}

					return rows, topics
				},
				initial: func() ([][]string, []*core.Topic) {
					rows := make([][]string, 0)
					topics := getApi().Core.GetTopics()
					for _, t := range topics {
						rows = append(rows, []string{
							t.Name,
						})
					}

					return rows, topics
				},
				onSelect: func(a *core.Topic) {
					if err := getApi().Core.MoveWorkspace(curWorkspace, a); err != nil {
						openToastDialogError(err.Error())
						return
					}

					if *sd != nil {
						(*sd).close()
					}

					getTopicsView().tableRenderer.SelectRow(0)

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
		Set('D', "Delete a workspace", func() {
			if getApi().Core.GetWorkspacesByTopicCount(getTopicsView().getSelectedTopic()) <= 0 {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					curWorkspace := wv.getSelectedWorkspace()
					getApi().Core.DeleteWorkspace(curWorkspace)

					// HACK: same as below
					getTopicsView().tableRenderer.SelectRow(0)
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		Set('r', "Rename workspace", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			openEditorDialogWithDefaultValue(func(s string) {
				if err := getApi().Core.RenameWorkspace(curWorkspace, s); err != nil {
					openToastDialogError(err.Error())
					return
				}
			}, func() {}, "New workspace name", smallEditorSize, curWorkspace.Name)
		}).
		Set('e', "Add/change description", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			openEditorDialog(func(desc string) {
				if desc != "" {
					getApi().Core.SetDescription(desc, curWorkspace)
				}
			}, func() {}, "Description", largeEditorSize)
		}).
		Set('a', "Create a workspace", func() {
			curTopic := getTopicsView().getSelectedTopic()
			if curTopic == nil {
				openToastDialog("You must create a topic first", false, "Note", func() {})
				return
			}

			openEditorDialog(func(name string) {
				if _, err := getApi().Core.CreateWorkspace(name, curTopic); err != nil {
					openToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new workspace
				// This will result in the workspace and the corresponding topic going to the top
				// because we are sorting by modifed time
				getTopicsView().tableRenderer.SelectRow(0)
				wv.tableRenderer.SelectRow(0)
			}, func() {}, "Workspace name ", smallEditorSize)
		}).
		Set('X', "Kill tmux session", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if getApi().Tmux.GetTmuxSessionByName(curWorkspace.Path) != nil {
				openConfirmationDialog(func(b bool) {
					if b {
						getApi().Core.DeleteWorkspaceTmuxSession(curWorkspace)
					}
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(wv.view.GetKeybindings(), func() {})
		})
}

func (wv *workspacesView) selectWorkspaceByShortPath(shortPath string) {
	wv.tableRenderer.SelectRowByValue(func(w *core.Workspace) bool {
		return w.ShortPath() == shortPath
	})
}

func (wv *workspacesView) refresh() {
	var workspaces core.Workspaces
	if selectedTopic := getTopicsView().getSelectedTopic(); selectedTopic != nil {
		workspaces = getApi().Core.GetWorkspaces().FilterByTopic(selectedTopic)
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
			if tm := getApi().Tmux.GetTmuxSessionByName(w.Path); tm != nil {
				numWindows := strconv.Itoa(tm.Windows)
				return numWindows + " - tmux"
			}

			return ""
		}()

		remote, err := w.GetGitRemote()
		if err != nil {
			openToastDialogError(err.Error())
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

func (wv *workspacesView) getSelectedWorkspace() *core.Workspace {
	_, w := wv.tableRenderer.GetSelectedRow()
	if w != nil {
		return *w
	}

	return nil
}

func (wv *workspacesView) render() error {
	wv.view.Clear()

	isFocused := wv.view.IsFocused()

	wv.tableRenderer.RenderWithSelectCallBack(wv.view, func(_ int, _ *tui.TableRow[*core.Workspace]) bool {
		return isFocused
	})

	return nil
}
