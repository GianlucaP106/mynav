package ui

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/utils"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type WorkspacesState struct {
	listRenderer *ListRenderer
	viewName     string
	workspaces   core.Workspaces
}

func newWorkspacesState() *WorkspacesState {
	workspace := &WorkspacesState{
		viewName: "WorkspacesView",
	}
	return workspace
}

func (ui *UI) initWorkspacesView() *gocui.View {
	view := ui.getView(ui.workspaces.viewName)
	if view != nil {
		return view
	}
	view = ui.setView(ui.workspaces.viewName)

	view.FrameColor = gocui.ColorBlue
	// view.FrameRunes = ThickFrame
	view.Title = withSurroundingSpaces("Workspaces")
	view.TitleColor = gocui.ColorBlue

	_, sizeY := view.Size()
	ui.workspaces.listRenderer = newListRenderer(0, sizeY/3, 0)

	ui.keyBinding(ui.workspaces.viewName).
		set('j', func() {
			ui.workspaces.listRenderer.increment()
		}).
		set('k', func() {
			ui.workspaces.listRenderer.decrement()
		}).
		set(gocui.KeyEsc, func() {
			ui.setFocusedFsView(ui.topics.viewName)
		}).
		set(gocui.KeyEnter, func() {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			if err := curWorkspace.OpenWorkspace(); err != nil {
				ui.openToastDialog(err.Error())
				return
			}
		}).
		set('d', func() {
			if ui.controller.GetWorkspacesByTopicCount(ui.getSelectedTopic()) <= 0 {
				return
			}

			ui.openConfirmationDialog(func(b bool) {
				if b {
					curWorkspace := ui.getSelectedWorkspace()
					ui.controller.DeleteWorkspace(curWorkspace)
					// HACK:
					ui.topics.listRenderer.setSelected(0)
					ui.refreshWorkspaces()
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		set('a', func() {
			curTopic := ui.getSelectedTopic()
			ui.openEditorDialog(func(name string) {
				ui.openEditorDialog(func(repoUrl string) {
					if err := ui.controller.CreateWorkspace(name, repoUrl, curTopic); err != nil {
						ui.openToastDialog(err.Error())
						return
					}

					// HACK: when there a is a new workspace
					// This will result in the workspace and the corresponding topic going to the top
					// because we are sorting by modifed time
					ui.topics.listRenderer.setSelected(0)
					ui.workspaces.listRenderer.setSelected(0)
					ui.refreshWorkspaces()
				}, func() {}, "Repo URL (leave blank if none)")
			}, func() {}, "Workspace name ")
		})
	return view
}

func (ui *UI) getSelectedWorkspace() *core.Workspace {
	return ui.getDisplayedWorkspace(ui.workspaces.listRenderer.selected)
}

func (ui *UI) getDisplayedWorkspace(idx int) *core.Workspace {
	wv := ui.workspaces.workspaces
	if idx >= len(wv) || idx < 0 {
		return nil
	}
	return wv[idx]
}

func (ui *UI) refreshWorkspaces() {
	ui.refreshTopics()
	out := ui.controller.GetWorkspacesByTopic(ui.getSelectedTopic())
	ui.workspaces.workspaces = out

	if ui.workspaces.listRenderer != nil {
		newListSize := len(ui.workspaces.workspaces)
		if ui.workspaces.listRenderer.listSize != newListSize {
			ui.workspaces.listRenderer.setListSize(newListSize)
		}
	}
}

func (ui *UI) formatWorkspaceItem(workspace *core.Workspace, selected bool) []string {
	sizeX, _ := ui.getView(ui.workspaces.viewName).Size()
	style, blankLine := func() (color.Style, string) {
		if selected {
			return color.New(color.Black, color.BgCyan), highlightedBlankLine(sizeX)
		}
		return color.New(color.Blue), blankLine(sizeX)
	}()

	lastModTime := workspace.GetLastModifiedTimeFormatted()
	gitRemote := func() string {
		remote := workspace.GetGitRemote()
		if remote != "" {
			return utils.TrimGithubUrl(remote)
		}
		return ""
	}()

	name := withSpacePadding(workspace.Name, sizeX/5)
	url := withSpacePadding(gitRemote, (sizeX*2)/5)
	time := withSpacePadding(lastModTime, sizeX/5)

	return []string{
		blankLine,
		displayLine(name+url+time, Left, sizeX, style),
		blankLine,
	}
}

func (ui *UI) renderWorkspacesView() {
	view := ui.initWorkspacesView()

	view.Clear()
	content := func() []string {
		if ui.workspaces.workspaces == nil {
			return []string{}
		}
		ui.refreshWorkspaces()
		out := make([]string, 0)
		ui.workspaces.listRenderer.forEach(func(i int) {
			selected := (ui.fs.focusedTab == ui.workspaces.viewName) && (i == ui.workspaces.listRenderer.selected)

			// TODO: use go routines here to optimize (git remote takes a long time)
			workspace := ui.formatWorkspaceItem(ui.workspaces.workspaces[i], selected)
			out = append(out, workspace...)
		})

		return out
	}()
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
}
