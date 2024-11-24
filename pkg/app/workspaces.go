package app

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

// Workspaces view displaying the workspaces of the current topic.
type Workspaces struct {
	// internal view
	view *tui.View

	// table renderer
	table *tui.TableRenderer[*core.Workspace]

	// loading flag to display loading
	loading bool
}

func newWorkspcacesView() *Workspaces {
	w := &Workspaces{}
	return w
}

func (wv *Workspaces) selectWorkspace(w *core.Workspace) {
	wv.table.SelectRowByValue(func(w2 *core.Workspace) bool {
		return w.ShortPath() == w2.ShortPath()
	})
}

func (wv *Workspaces) selected() *core.Workspace {
	_, w := wv.table.SelectedRow()
	if w != nil {
		return *w
	}

	return nil
}

func (w *Workspaces) setLoading(l bool) {
	w.loading = l
}

func (w *Workspaces) getLoading() bool {
	return w.loading
}

func (wv *Workspaces) focus() {
	a.focusView(wv.view)
	wv.show()
}

func (wv *Workspaces) show() {
	def := "No preview"
	w := wv.selected()
	if w == nil {
		a.preview.show(def)
		return
	}

	a.info.show(w)

	s := a.api.WorkspacePreview(w)
	if s == "" {
		a.preview.show(def)
		return
	}

	a.preview.show(s)
}

func (wv *Workspaces) refresh() {
	var workspaces core.Workspaces
	if selectedTopic := a.topics.selected(); selectedTopic != nil {
		workspaces = a.api.AllWorkspaces().ByTopic(selectedTopic)
	} else {
		workspaces = make(core.Workspaces, 0)
	}

	sMap := a.api.SessionMap()
	rows := make([][]string, 0)
	for _, w := range workspaces.Sorted() {
		tmux := "-"
		s := sMap.Get(w)
		if s != nil {
			tmux = "Yes"
		}
		rows = append(rows, []string{
			w.Name,
			tmux,
			w.LastModifiedTimeFormatted(),
		})
	}

	wv.table.Fill(rows, workspaces)
}

func (wv *Workspaces) render() {
	wv.view.Clear()
	a.ui.Resize(wv.view, getViewPosition(wv.view.Name()))

	// update title based on which topic is selected
	if t := a.topics.selected(); t != nil {
		wv.view.Title = fmt.Sprintf(" Workspaces - %s ", a.topics.selected().Name)
	} else {
		wv.view.Title = " Workspaces "
	}

	// update page row marker
	row, _ := wv.table.SelectedRow()
	size := wv.table.Size()
	wv.view.Subtitle = fmt.Sprintf(" %d / %d ", min(row+1, size), size)

	if wv.getLoading() {
		fmt.Fprintln(wv.view, "Loading...")
		return
	}

	isFocused := a.ui.IsFocused(wv.view)
	wv.table.RenderSelect(wv.view, func(_ int, _ *tui.TableRow[*core.Workspace]) bool {
		return isFocused
	})
}

func (wv *Workspaces) init() {
	wv.view = a.ui.SetView(getViewPosition(WorkspacesView))
	a.styleView(wv.view)

	sizeX, sizeY := wv.view.Size()
	titles := []string{
		"Name",
		"Session",
		"Last Modified",
	}
	proportions := []float64{
		0.40,
		0.20,
		0.40,
	}
	styles := []color.Style{
		color.New(color.FgBlue, color.Bold),
		color.Success.Style,
		color.New(color.FgDarkGray, color.OpItalic),
	}
	wv.table = tui.NewTableRenderer[*core.Workspace]()
	wv.table.Init(sizeX, sizeY, titles, proportions)
	wv.table.SetStyles(styles)

	down := func() {
		wv.table.Down()
		wv.show()
	}

	up := func() {
		wv.table.Up()
		wv.show()
	}
	tv := a.topics
	a.ui.KeyBinding(wv.view).
		Set('j', "Move down", down).
		Set('k', "Move up", up).
		Set(gocui.KeyArrowDown, "Move down", down).
		Set(gocui.KeyArrowUp, "Move up", up).
		Set(gocui.KeyEsc, "Go back", func() {
			tv.focus()
		}).
		Set('c', "Command", func() {
			w := wv.selected()
			if w == nil {
				return
			}

			e := editor(func(s string) {
				split := strings.Split(s, " ")
				if len(split) == 0 {
					toast("command is empty", toastError)
					return
				}

				if !system.DoesProgramExist(split[0]) {
					toast(fmt.Sprintf("%s is not installed on the system", s), toastError)
					return
				}

				split = append(split, w.Path())

				var err error
				a.runAction(func() {
					err = system.CommandWithRedirect(split...).Run()
				})
				if err != nil {
					toast(err.Error(), toastError)
				}

				a.api.SelectWorkspace(w)
			}, func() {}, "Command", smallEditorSize, "nvim")
			e.view.Subtitle = " Workspace path will be appended "
		}).
		Set('g', "Clone git repo", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			editor(func(s string) {
				if err := a.api.CloneWorkspaceRepo(curWorkspace, s); err != nil {
					toast(err.Error(), toastError)
				}

				a.refreshAll()
				toast("Cloned repo to workspace "+curWorkspace.Name, toastInfo)
			}, func() {}, "Git repo URL", smallEditorSize, "")
		}).
		Set('G', "Open browser to git repo", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			remote, err := curWorkspace.GitRemote()
			if err != nil {
				return
			}

			if err := system.OpenBrowser(remote); err != nil {
				toast(err.Error(), toastError)
			}
		}).
		Set('u', "Copy git repo url to clipboard", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			remote, err := curWorkspace.GitRemote()
			if err != nil {
				return
			}

			system.CopyToClip(remote)
			toast("Copied "+remote+" to clipboard", toastInfo)
		}).
		Set(gocui.KeyEnter, "Open workspace", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			var error error
			a.runAction(func() {
				error = a.api.OpenWorkspace(curWorkspace)
			})
			if error != nil {
				toast(error.Error(), toastError)
			}
			a.refreshAll()
		}).
		Set('m', "Move workspace", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			sd := new(*Search[*core.Topic])
			*sd = search(searchDialogConfig[*core.Topic]{
				onSearch: func(s string) ([][]string, []*core.Topic) {
					rows := make([][]string, 0)
					topics := a.api.AllTopics().ByNameContaining(s)
					for _, t := range topics {
						rows = append(rows, []string{
							t.Name,
						})
					}

					return rows, topics
				},
				initial: func() ([][]string, []*core.Topic) {
					rows := make([][]string, 0)
					topics := a.api.AllTopics()
					for _, t := range topics {
						rows = append(rows, []string{
							t.Name,
						})
					}

					return rows, topics
				},
				onSelect: func(t *core.Topic) {
					if err := a.api.MoveWorkspace(curWorkspace, t); err != nil {
						toast(err.Error(), toastError)
						return
					}

					if *sd != nil {
						(*sd).close()
					}

					a.refresh(curWorkspace.Topic, nil, true, false)
					toast("Moved workspace "+curWorkspace.Name, toastInfo)
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
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			alert(func(b bool) {
				if b {
					t := curWorkspace.Topic
					a.api.DeleteWorkspace(curWorkspace)
					a.refresh(t, nil, true, false)
					toast("Deleted workspace "+curWorkspace.Name, toastInfo)
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		Set('r', "Rename workspace", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			editor(func(s string) {
				if err := a.api.RenameWorkspace(curWorkspace, s); err != nil {
					toast(err.Error(), toastError)
					return
				}

				a.refresh(curWorkspace.Topic, curWorkspace, true, false)
				toast("Renamed workspace "+curWorkspace.Name, toastInfo)
			}, func() {}, "New workspace name", smallEditorSize, curWorkspace.Name)
		}).
		Set('a', "Create a workspace", func() {
			curTopic := a.topics.selected()
			if curTopic == nil {
				toast("You must create a topic first", toastWarn)
				return
			}

			editor(func(name string) {
				w, err := a.api.NewWorkspace(curTopic, name)
				if err != nil {
					toast(err.Error(), toastError)
					return
				}

				a.refresh(curTopic, w, true, false)
				toast("Created workspace "+w.Name, toastInfo)
			}, func() {}, "Workspace name ", smallEditorSize, "")
		}).
		Set('X', "Kill session", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			if a.api.Session(curWorkspace) != nil {
				alert(func(b bool) {
					if !b {
						return
					}

					if err := a.api.KillSession(curWorkspace); err != nil {
						toast(err.Error(), toastError)
						return
					}

					a.refreshAll()
					toast("Killed session "+curWorkspace.Name, toastInfo)
				}, "Are you sure you want to delete the session?")
			}
		}).
		Set('h', "Focus topics view", func() {
			a.topics.focus()
		}).
		Set(gocui.KeyArrowLeft, "Focus topics view", func() {
			a.topics.focus()
		}).
		Set('l', "Focus sessions view", func() {
			a.sessions.focus()
		}).
		Set(gocui.KeyArrowRight, "Focus sessions view", func() {
			a.sessions.focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			help(wv.view)
		})
}
