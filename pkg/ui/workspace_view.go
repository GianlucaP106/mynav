package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"os/exec"
	"strconv"
	"strings"

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

func (wv *workspacesView) focus() {
	focusView(wv.getView().Name())
}

func (wv *workspacesView) init() {
	wv.view = getViewPosition(WorkspacesView).Set()

	wv.view.Title = tui.WithSurroundingSpaces("Workspaces")
	styleView(wv.getView())

	sizeX, sizeY := wv.view.Size()

	titles := []string{
		"Name",
		"Description",
		"Git Remote",
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

	wv.refresh()

	if selectedWorkspace := getApi().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		wv.selectWorkspaceByShortPath(selectedWorkspace.ShortPath())
	}

	displayedWorkspaceOpenerCmd := getApi().GlobalConfiguration.GetCustomWorkspaceOpenerCmd()
	displayedWorkspaceOpenerCmdStr := "tmux/nvim"
	if len(displayedWorkspaceOpenerCmd) > 0 {
		displayedWorkspaceOpenerCmdStr = strings.Join(displayedWorkspaceOpenerCmd, " ")
	}

	tv := getTopicsView()
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

			tv.focus()
		}).
		Set('s', "See workspace information", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			openWorkspaceInfoDialog(curWorkspace, func() {})
		}).
		Set('L', "Open lazygit at this workspace (if there is a git repository)", func() {
			w := wv.getSelectedWorkspace()
			if w == nil {
				return
			}

			if !system.DoesProgramExist("lazygit") {
				openToastDialogError("lazygit is not installed on the system")
				return
			}

			if w.GitRemote == nil {
				return
			}

			runAction(func() {
				system.OpenLazygit(w.Path)
			})
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

				refreshMainViews()
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
		Set('u', "Copy git repo url to clipboard", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if curWorkspace.GitRemote == nil {
				return
			}
			remote := *curWorkspace.GitRemote

			system.CopyToClip(remote)
			openToastDialog(remote, toastDialogNeutralType, "Repo URL copied", func() {})
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
		Set(gocui.KeyEnter, "Open in "+displayedWorkspaceOpenerCmdStr, func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			var error error = nil
			runAction(func() {
				error = getApi().Core.OpenWorkspace(curWorkspace)
			})

			if error != nil {
				openToastDialogError(error.Error())
			}
		}).
		Set('v', "Open in neovim", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			runAction(func() {
				getApi().Core.OpenNeovimInWorkspace(curWorkspace)
			})
		}).
		Set('t', "Open in terminal", func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			terminalOpenCmd := getApi().GlobalConfiguration.GetTerminalOpenerCmd()
			var error error = nil
			runAction(func() {
				if len(terminalOpenCmd) > 0 {
					err := exec.Command(terminalOpenCmd[0], terminalOpenCmd[1:]...).Run()
					if err != nil {
						error = err
					}
				} else {
					getApi().Core.OpenTerminalInWorkspace(curWorkspace)
				}
			})

			openToastDialogError(error.Error())
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

					tv.tableRenderer.SelectRow(0)
					refreshMainViews()
					wv.focus()
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
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					getApi().Core.DeleteWorkspace(curWorkspace)

					// HACK: same as below
					tv.tableRenderer.SelectRow(0)
					refreshMainViews()
					refreshTmuxViews()
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

				tv.tableRenderer.SelectRow(0)
				wv.tableRenderer.SelectRow(0)
				refreshTmuxViews()
				refreshMainViews()
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
					refreshMainViews()
				}
			}, func() {}, "Description", largeEditorSize)
		}).
		Set('a', "Create a workspace", func() {
			curTopic := getTopicsView().getSelectedTopic()
			if curTopic == nil {
				openToastDialog("You must create a topic first", toastDialogNeutralType, "Note", func() {})
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
				tv.tableRenderer.SelectRow(0)
				wv.tableRenderer.SelectRow(0)
				refreshMainViews()
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
						refreshMainViews()
						refreshTmuxViews()
					}
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(wv.view.GetKeybindings(), func() {})
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
			remote = core.TrimGithubUrl(remote)
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
	isFocused := wv.view.IsFocused()
	wv.view.Clear()
	wv.view.Resize(getViewPosition(wv.view.Name()))

	wv.tableRenderer.RenderWithSelectCallBack(wv.view, func(_ int, _ *tui.TableRow[*core.Workspace]) bool {
		return isFocused
	})

	return nil
}
