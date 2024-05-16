package ui

import (
	"fmt"
	"mynav/pkg/api"
	"mynav/pkg/utils"
	"strconv"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type WorkspacesState struct {
	listRenderer *ListRenderer
	viewName     string
	search       string
	workspaces   api.Workspaces
}

func newWorkspacesState() *WorkspacesState {
	workspace := &WorkspacesState{
		viewName: "WorkspacesView",
	}
	return workspace
}

func (ui *UI) initWorkspacesView() *gocui.View {
	exists := false
	view := ui.getView(ui.workspaces.viewName)
	exists = view != nil
	if !exists {
		view = ui.setView(ui.workspaces.viewName)
	}

	if ui.workspaces.search != "" {
		view.Subtitle = withSurroundingSpaces("Searching: " + ui.workspaces.search)
	} else {
		view.Subtitle = ""
	}

	view.Title = withSurroundingSpaces("Workspaces")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorBlue

	if exists {
		return view
	}

	_, sizeY := view.Size()
	ui.workspaces.listRenderer = newListRenderer(0, sizeY/3, 0)
	ui.refreshWorkspaces()

	if selectedWorkspace := ui.controller.WorkspaceManager.GetSelectedWorkspace(); selectedWorkspace != nil {
		ui.selectWorkspaceByShortPath(selectedWorkspace.ShortPath())
	}

	ui.keyBinding(ui.workspaces.viewName).
		set('j', func() {
			ui.workspaces.listRenderer.increment()
		}).
		set('k', func() {
			ui.workspaces.listRenderer.decrement()
		}).
		set(gocui.KeyEsc, func() {
			if ui.workspaces.search != "" {
				ui.workspaces.search = ""
				return
			}

			ui.setFocusedFsView(ui.topics.viewName)
		}).
		set('s', func() {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			ui.openWorkspaceInfoDialog(curWorkspace)
		}).
		set('/', func() {
			ui.openEditorDialog(func(s string) {
				ui.workspaces.search = s
			}, func() {}, "Search", Small)
		}).
		setKeybinding(ui.workspaces.viewName, gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			if utils.IsTmuxSession() {
				ui.setWorkspaceAction(curWorkspace, utils.NvimCmd(curWorkspace.Path))
				return gocui.ErrQuit
			}

			foundExisting, sessionName := ui.controller.WorkspaceManager.GetOrCreateTmuxSession(curWorkspace)
			if foundExisting {
				ui.setWorkspaceAction(curWorkspace, utils.AttachTmuxSessionCmd(sessionName))
			} else {
				ui.setWorkspaceAction(curWorkspace, utils.NewTmuxSessionCmd(sessionName, curWorkspace.Path))
			}

			return gocui.ErrQuit
		}).
		setKeybinding(ui.workspaces.viewName, 'v', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			ui.setWorkspaceAction(curWorkspace, utils.NvimCmd(curWorkspace.Path))
			return gocui.ErrQuit
		}).
		setKeybinding(ui.workspaces.viewName, 't', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			openTermCmd, err := utils.GetOpenTerminalCmd(curWorkspace.Path)
			if err != nil {
				ui.openToastDialog(err.Error())
				return nil
			}

			ui.setWorkspaceAction(curWorkspace, openTermCmd)

			return gocui.ErrQuit
		}).
		set('d', func() {
			if ui.controller.
				WorkspaceManager.
				Workspaces.
				ByTopic(ui.getSelectedTopic()).Len() <= 0 {
				return
			}

			ui.openConfirmationDialog(func(b bool) {
				if b {
					curWorkspace := ui.getSelectedWorkspace()
					ui.controller.WorkspaceManager.DeleteWorkspace(curWorkspace)
					// HACK: same as below
					ui.topics.listRenderer.setSelected(0)
					ui.refreshTopics()
					ui.refreshWorkspaces()
				}
			}, "Are you sure you want to delete this workspace?")
		}).
		set('r', func() {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			ui.openEditorDialog(func(desc string) {
				if desc != "" {
					ui.controller.WorkspaceManager.SetDescription(curWorkspace, desc)
				}
			}, func() {}, "Description", Large)
		}).
		set('a', func() {
			curTopic := ui.getSelectedTopic()
			ui.openEditorDialog(func(name string) {
				ui.openEditorDialog(func(repoUrl string) {
					if _, err := ui.controller.WorkspaceManager.CreateWorkspace(name, repoUrl, curTopic); err != nil {
						ui.openToastDialog(err.Error())
						return
					}

					// HACK: when there a is a new workspace
					// This will result in the workspace and the corresponding topic going to the top
					// because we are sorting by modifed time
					ui.topics.listRenderer.setSelected(0)
					ui.workspaces.listRenderer.setSelected(0)
					ui.refreshTopics()
					ui.refreshWorkspaces()
				}, func() {}, "Repo URL (leave blank if none)", Small)
			}, func() {}, "Workspace name ", Small)
		})
	return view
}

func (ui *UI) setWorkspaceAction(w *api.Workspace, action []string) {
	ui.setAction(action)
	ui.controller.WorkspaceManager.WorkspaceStore.SetSelectedWorkspace(w)
}

func (ui *UI) getSelectedWorkspace() *api.Workspace {
	return ui.getDisplayedWorkspace(ui.workspaces.listRenderer.selected)
}

func (ui *UI) getDisplayedWorkspace(idx int) *api.Workspace {
	wv := ui.workspaces.workspaces
	if idx >= len(wv) || idx < 0 {
		return nil
	}
	return wv[idx]
}

func (ui *UI) refreshWorkspaces() {
	workspaces := ui.controller.WorkspaceManager.Workspaces.ByTopic(ui.getSelectedTopic())

	if ui.workspaces.search != "" {
		workspaces = workspaces.FilterByNameContaining(ui.workspaces.search)
	}

	ui.workspaces.workspaces = workspaces

	if ui.workspaces.listRenderer != nil {
		newListSize := len(ui.workspaces.workspaces)
		if ui.workspaces.listRenderer.listSize != newListSize {
			ui.workspaces.listRenderer.setListSize(newListSize)
		}
	}
}

func (ui *UI) formatWorkspaceRow(workspace *api.Workspace, selected bool) []string {
	sizeX, _ := ui.getView(ui.workspaces.viewName).Size()
	style, blank := func() (color.Style, string) {
		if selected {
			return color.New(color.Black, color.BgCyan), highlightedBlankLine(sizeX + 5) // +5 for extra padding
		}
		return color.New(color.Blue), blankLine(sizeX)
	}()

	lastModTime := workspace.GetLastModifiedTimeFormatted()
	gitRemote := func() string {
		remote := ui.controller.WorkspaceManager.GetGitRemote(workspace)
		if remote != "" {
			return utils.TrimGithubUrl(remote)
		}
		return ""
	}()

	fifth := sizeX / 5
	description := withSpacePadding(workspace.GetDescription(), fifth)
	url := withSpacePadding(gitRemote, fifth)
	time := withSpacePadding(lastModTime, fifth)

	name := withSurroundingSpaces(workspace.Name)
	tmux := func() string {
		if workspace.Metadata.TmuxSession != nil {
			tm := workspace.Metadata.TmuxSession
			numWindows := strconv.Itoa(tm.NumWindows)
			var line string
			if selected {
				c := color.New(color.BgCyan, color.Black)
				numWindows = c.Sprint(numWindows)
				line = numWindows + c.Sprint(" - tmux")
			} else {
				numWindows = color.New(color.BgGreen, color.Black).Sprint(numWindows)
				line = numWindows + color.Green.Sprint(" - tmux")
			}
			return line
		}
		return ""
	}()

	name = withSpacePadding(name, sizeX/5)
	line := style.Sprint(name+description+url+time) + tmux + style.Sprint(strings.Repeat(" ", fifth+5)) // +5 for extra padding
	return []string{
		blank,
		line,
		blank,
	}
}

func (ui *UI) selectWorkspaceByShortPath(shortPath string) {
	for idx, w := range ui.workspaces.workspaces {
		if w.ShortPath() == shortPath {
			ui.workspaces.listRenderer.setSelected(idx)
		}
	}
}

func (ui *UI) renderWorkspacesView() {
	view := ui.initWorkspacesView()

	view.Clear()
	content := func() []string {
		if ui.workspaces.workspaces == nil {
			return []string{}
		}
		out := make([]string, 0)
		ui.workspaces.listRenderer.forEach(func(i int) {
			selected := (ui.fs.focusedTab == ui.workspaces.viewName) && (i == ui.workspaces.listRenderer.selected)

			// TODO: https://github.com/GianlucaP106/mynav/issues/18
			workspace := ui.formatWorkspaceRow(ui.workspaces.workspaces[i], selected)
			out = append(out, workspace...)
		})

		return out
	}()
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
}
