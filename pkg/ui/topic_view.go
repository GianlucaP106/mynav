package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type topicsView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*core.Topic]
	search        *core.Value[string]
	globalSearch  *core.Value[string]
}

var _ viewable = new(topicsView)

func newTopicsView() *topicsView {
	return &topicsView{
		search: core.NewValue(""),
	}
}

func (tv *topicsView) getView() *tui.View {
	return tv.view
}

func (tv *topicsView) focus() {
	ui.focusView(tv.getView().Name())
}

func (tv *topicsView) init() {
	tv.view = getViewPosition(TopicView).Set()

	tv.view.Title = tui.WithSurroundingSpaces("Topics")
	ui.styleView(tv.view)

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = tui.NewTableRenderer[*core.Topic]()
	titles := []string{
		"Name",
		"Last Modified",
	}
	colProportions := []float64{
		0.5,
		0.5,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, colProportions)

	tv.refresh()

	if selectedWorkspace := api().Workspaces.GetSelectedWorkspace(); selectedWorkspace != nil {
		tv.selectTopicByName(selectedWorkspace.Topic.Name)
	}

	moveRight := func() {
		if api().Topics.GetTopicCount() > 0 {
			ui.getWorkspacesView().focus()
		}
	}

	wv := ui.getWorkspacesView()
	tv.view.KeyBinding().
		Set('j', "Move down", func() {
			tv.tableRenderer.Down()
			ui.refresh(wv)
		}).
		Set('k', "Move up", func() {
			tv.tableRenderer.Up()
			ui.refresh(wv)
		}).
		Set(gocui.KeyEnter, "Open topic", moveRight).
		Set('/', "Search by name", func() {
			openEditorDialog(func(s string) {
				tv.search.Set(s)
				tv.view.Subtitle = tui.WithSurroundingSpaces("Searching: " + s)
				refreshMainViews()
			}, func() {}, "Search", smallEditorSize)
		}).
		Set(gocui.KeyEsc, "Escape search", func() {
			if tv.search.Get() != "" {
				tv.search.Set("")
				tv.view.Subtitle = ""
				refreshMainViews()
			}
		}).
		Set('a', "Create a topic", func() {
			openEditorDialog(func(s string) {
				if err := api().Topics.CreateTopic(s); err != nil {
					openToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.tableRenderer.SelectRow(0)
				refreshMainViews()
			}, func() {}, "Topic name", smallEditorSize)
		}).
		Set('r', "Rename topic", func() {
			t := tv.getSelectedTopic()
			if t == nil {
				return
			}

			openEditorDialogWithDefaultValue(func(s string) {
				if err := api().Topics.RenameTopic(t, s); err != nil {
					openToastDialogError(err.Error())
					return
				}

				refreshMainViews()
			}, func() {}, "New topic name", smallEditorSize, t.Name)
		}).
		Set('s', "Search for a workspace", func() {
			sd := new(*searchListDialog[*core.Workspace])
			*sd = openSearchListDialog(searchDialogConfig[*core.Workspace]{
				onSearch: func(s string) ([][]string, []*core.Workspace) {
					workspaces := api().Workspaces.GetWorkspaces().Sorted().FilterByNameContaining(s)
					rows := make([][]string, 0)
					for _, w := range workspaces {
						remote, _ := w.GetGitRemote()
						if remote != "" {
							remote = core.TrimGithubUrl(remote)
						}

						session := api().Tmux.GetTmuxSessionByName(w.Path)
						sessionStr := "None"
						if session != nil {
							sessionStr = "Active"
						}

						rows = append(rows, []string{
							w.ShortPath(),
							remote,
							sessionStr,
						})
					}

					return rows, workspaces
				},
				onSelect: func(w *core.Workspace) {
					tv.selectTopicByName(w.Topic.Name)
					wv := ui.getWorkspacesView()
					wv.refresh()
					wv.selectWorkspaceByShortPath(w.ShortPath())

					if *sd != nil {
						(*sd).close()
					}

					wv.focus()
				},
				onSelectDescription: "Go to workspace",
				searchViewTitle:     "Search a workspace",
				tableViewTitle:      "Result",
				tableTitles: []string{
					"Workspace",
					"Git Remote",
					"Tmux Session",
				}, tableProportions: []float64{
					0.4,
					0.4,
					0.2,
				},
			})
		}).
		Set('D', "Delete topic", func() {
			if api().Topics.GetTopicCount() <= 0 {
				return
			}
			openConfirmationDialog(func(b bool) {
				if !b {
					return
				}

				if err := api().Topics.DeleteTopic(tv.getSelectedTopic()); err != nil {
					openToastDialogError(err.Error())
				}

				refreshMainViews()
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(tv.view.GetKeybindings(), func() {})
		})
}

func (tv *topicsView) refresh() {
	topics := api().Topics.GetTopics().Sorted()

	search := tv.search.Get()
	if search != "" {
		topics = topics.FilterByNameContaining(search)
	}

	rows := make([][]string, 0)
	rowValues := make([]*core.Topic, 0)
	for _, topic := range topics {
		rowValues = append(rowValues, topic)
		rows = append(rows, []string{
			topic.Name,
			topic.GetLastModifiedTimeFormatted(),
		})
	}

	tv.tableRenderer.FillTable(rows, rowValues)
}

func refreshMainViews() {
	if !api().GlobalConfiguration.Standalone {
		ui.queueRefresh(func() {
			t := ui.getTopicsView()
			t.refresh()
			ui.renderView(t)

			wv := ui.getWorkspacesView()
			wv.refresh()
			ui.renderView(wv)
		})
	}
}

func (tv *topicsView) getSelectedTopic() *core.Topic {
	_, t := tv.tableRenderer.GetSelectedRow()
	if t != nil {
		return *t
	}
	return nil
}

func (tv *topicsView) selectTopicByName(name string) {
	tv.tableRenderer.SelectRowByValue(func(t *core.Topic) bool {
		return t.Name == name
	})
}

func (tv *topicsView) render() error {
	currentViewSelected := tv.view.IsFocused()
	tv.view.Clear()
	tv.view.Resize(getViewPosition(tv.view.Name()))

	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *tui.TableRow[*core.Topic]) bool {
		return currentViewSelected
	})

	return nil
}
