package app

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/GianlucaP106/mynav/pkg/core"
	"github.com/GianlucaP106/mynav/pkg/tui"
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
		return w.Path() == w2.Path()
	})
}

func (wv *Workspaces) selected() *core.Workspace {
	_, w := wv.table.SelectedRow()
	if w != nil {
		return w.Value
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
	wv.refreshDown()
}

func (wv *Workspaces) showInfo() {
	w := wv.selected()
	if w == nil {
		a.info.show(nil)
		return
	}

	a.info.show(w)
}

func (wv *Workspaces) refreshPreview() {
	w := wv.selected()
	if w == nil {
		a.preview.setSession(nil)
		return
	}

	s := a.api.Session(w)
	a.preview.setSession(s)
}

func (wv *Workspaces) refreshDown() {
	wv.showInfo()
	a.worker.Queue(func() {
		wv.refreshPreview()
		a.ui.Update(func() {
			a.preview.render()
		})
	})
}

func (wv *Workspaces) refresh() {
	var workspaces core.Workspaces
	if selectedTopic := a.topics.selected(); selectedTopic != nil {
		workspaces = a.api.Workspaces(selectedTopic)
	} else {
		workspaces = make(core.Workspaces, 0)
	}

	sMap := a.api.SessionMap()
	tableRows := make([]*tui.TableRow[*core.Workspace], 0)
	for _, w := range workspaces.Sorted() {
		tmux := ""
		s := sMap.Get(w)
		if s != nil {
			tmux = "Yes"
		}
		timeStr := core.TimeAgo(w.LastModified())
		tableRows = append(tableRows, &tui.TableRow[*core.Workspace]{
			Cols: []string{
				w.Name,
				tmux,
				timeStr,
			},
			Value: w,
		})
	}

	wv.table.Fill(tableRows)
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
	wv.table.RenderTable(wv.view, func(_ int, _ *tui.TableRow[*core.Workspace]) bool {
		return isFocused
	}, func(i int, tr *tui.TableRow[*core.Workspace]) {
		newTime := core.TimeAgo(tr.Value.LastModified())
		tr.Cols[len(tr.Cols)-1] = newTime
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
		workspaceNameColor,
		sessionMarkerColor,
		timestampColor,
	}
	wv.table = tui.NewTableRenderer[*core.Workspace]()
	wv.table.Init(sizeX, sizeY, titles, proportions)
	wv.table.SetStyles(styles)

	down := func() {
		wv.table.Down()
		wv.refreshDown()
	}

	up := func() {
		wv.table.Up()
		wv.refreshDown()
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

				if !core.DoesProgramExist(split[0]) {
					toast(fmt.Sprintf("%s is not installed on the system", s), toastError)
					return
				}

				split = append(split, w.Path())

				err := a.runAction(func() error {
					return core.CommandWithRedirect(split...).Run()
				})
				if err != nil {
					toast(err.Error(), toastError)
				}

				a.api.SelectWorkspace(w)
			}, func() {}, "Command", smallEditorSize, "nvim")
			e.view.Subtitle = " Workspace path will be appended "
		}).
		Set('i', "Clone git repo", func() {
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
		Set('I', "Open browser to git repo", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			remote, err := curWorkspace.GitRemote()
			if err != nil {
				return
			}

			if err := core.OpenBrowser(remote); err != nil {
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

			core.CopyToClip(remote)
			toast("Copied "+remote+" to clipboard", toastInfo)
		}).
		Set(gocui.KeyEnter, "Open workspace", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			if core.IsTmuxSession() {
				toast("A tmux session is already active", toastWarn)
				return
			}

			start := time.Now()
			err := a.runAction(func() error {
				return a.api.OpenWorkspace(curWorkspace)
			})
			if err != nil {
				toast(err.Error(), toastError)
			} else {
				timeTaken := time.Since(start)
				s := fmt.Sprintf("Detached session %s - %s active", curWorkspace.Name, core.TimeDeltaStr(timeTaken))
				toast(s, toastInfo)
			}

			a.refresh(curWorkspace.Topic, curWorkspace, nil)
		}).
		Set('m', "Move workspace", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			sd := new(*Search[*core.Topic])
			*sd = search(SearchDialogConfig[*core.Topic]{
				onSearch: func(s string) []*tui.TableRow[*core.Topic] {
					rows := make([]*tui.TableRow[*core.Topic], 0)
					topics := a.api.Topics()
					filtered := core.Topics{}
					for _, topic := range topics {
						if strings.Contains(topic.Name, s) {
							filtered = append(filtered, topic)
						}
					}

					for _, t := range filtered {
						rows = append(rows, &tui.TableRow[*core.Topic]{
							Cols: []string{
								t.Name,
							},
							Value: t,
						})
					}

					return rows
				},
				initial: func() []*tui.TableRow[*core.Topic] {
					rows := make([]*tui.TableRow[*core.Topic], 0)
					topics := a.api.Topics()
					for _, t := range topics {
						rows = append(rows, &tui.TableRow[*core.Topic]{
							Cols: []string{
								t.Name,
							},
							Value: t,
						})
					}

					return rows
				},
				onSelect: func(t *core.Topic) {
					if err := a.api.MoveWorkspace(curWorkspace, t); err != nil {
						toast(err.Error(), toastError)
						return
					}

					if *sd != nil {
						(*sd).close()
					}

					a.refresh(curWorkspace.Topic, nil, nil)
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
					a.refresh(t, nil, nil)
					toast("Deleted workspace "+curWorkspace.Name, toastInfo)
				}
			}, fmt.Sprintf("Are you sure you want to delete workspace %s?", curWorkspace.Name))
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

				a.refresh(curWorkspace.Topic, curWorkspace, nil)
				toast("Renamed workspace "+curWorkspace.Name, toastInfo)
			}, func() {}, "New workspace name", smallEditorSize, curWorkspace.Name)
		}).
		Set('A', "Create a workspace from git url", func() {
			curTopic := a.topics.selected()
			if curTopic == nil {
				toast("You must create a topic first", toastWarn)
				return
			}

			editor(func(uri string) {
				_, err := url.ParseRequestURI(uri)
				if err != nil {
					toast("URL must be valid", toastError)
					return
				}

				urlTrimmed := strings.TrimSuffix(uri, ".git")
				name := path.Base(urlTrimmed)
				w, err := a.api.NewWorkspace(curTopic, name)
				if err != nil {
					toast(err.Error(), toastError)
					return
				}

				if err := a.api.CloneWorkspaceRepo(w, uri); err != nil {
					toast(err.Error(), toastError)
					return
				}

				a.refresh(curTopic, w, nil)
				toast("Cloned and created workspace "+w.Name, toastInfo)
			}, func() {}, "Git url", smallEditorSize, "")
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

				a.refresh(curTopic, w, nil)
				toast("Created workspace "+w.Name, toastInfo)
			}, func() {}, "Name", smallEditorSize, "")
		}).
		Set('X', "Kill session", func() {
			curWorkspace := wv.selected()
			if curWorkspace == nil {
				return
			}

			if s := a.api.Session(curWorkspace); s != nil {
				alert(func(b bool) {
					if !b {
						return
					}

					if err := s.Kill(); err != nil {
						toast(err.Error(), toastError)
						return
					}

					a.refreshAll()
					toast("Killed session "+curWorkspace.Name, toastInfo)
				}, fmt.Sprintf("Are you sure you want to kill session for %s?", curWorkspace.Name))
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
		Set('g', "Go to top", func() {
			wv.table.Top()
			wv.refreshDown()
		}).
		Set('G', "Go to bottom", func() {
			wv.table.Bottom()
			wv.refreshDown()
		}).
		Set(gocui.KeyArrowRight, "Focus sessions view", func() {
			a.sessions.focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			help(wv.view)
		})
}
