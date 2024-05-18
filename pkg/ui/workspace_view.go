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

	if selectedWorkspace := ui.api.GetSelectedWorkspace(); selectedWorkspace != nil {
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
				ui.refreshWorkspaces()
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
		set('g', func() {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			ui.openEditorDialog(func(s string) {
				if err := ui.api.CloneRepo(s, curWorkspace); err != nil {
					ui.openToastDialog(err.Error())
				}
			}, func() {}, "Git repo URL", Small)
		}).
		set('/', func() {
			ui.openEditorDialog(func(s string) {
				ui.workspaces.search = s
				ui.refreshWorkspaces()
			}, func() {}, "Search", Small)
		}).
		setKeybinding(ui.workspaces.viewName, gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			if utils.IsTmuxSession() {
				ui.setAction(utils.NvimCmd(curWorkspace.Path))
				return gocui.ErrQuit
			}

			command := ui.api.CreateOrAttachTmuxSession(curWorkspace)
			ui.setAction(command)

			return gocui.ErrQuit
		}).
		setKeybinding(ui.workspaces.viewName, 'v', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			ui.api.SetSelectedWorkspace(curWorkspace)
			ui.setAction(utils.NvimCmd(curWorkspace.Path))
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

			ui.setAction(openTermCmd)
			return gocui.ErrQuit
		}).
		set('d', func() {
			if ui.api.GetWorkspacesByTopicCount(ui.getSelectedTopic()) <= 0 {
				return
			}

			ui.openConfirmationDialog(func(b bool) {
				if b {
					curWorkspace := ui.getSelectedWorkspace()
					ui.api.DeleteWorkspace(curWorkspace)
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
					ui.api.SetDescription(desc, curWorkspace)
				}
			}, func() {}, "Description", Large)
		}).
		set('a', func() {
			curTopic := ui.getSelectedTopic()
			ui.openEditorDialog(func(name string) {
				if _, err := ui.api.CreateWorkspace(name, curTopic); err != nil {
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
			}, func() {}, "Workspace name ", Small)
		}).
		set('x', func() {
			curWorkspace := ui.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if curWorkspace.Metadata.TmuxSession != nil {
				ui.openConfirmationDialog(func(b bool) {
					if b {
						ui.api.DeleteTmuxSession(curWorkspace)
					}
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		set('?', func() {
			ui.openHelpView(ui.getKeyBindings(ui.workspaces.viewName))
		})
	return view
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
	var workspaces api.Workspaces
	if selectedTopic := ui.getSelectedTopic(); selectedTopic != nil {
		workspaces = ui.api.GetWorkspaces().ByTopic(ui.getSelectedTopic())
	} else {
		workspaces = make(api.Workspaces, 0)
	}

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
		// TODO: handle error
		remote, _ := workspace.GetGitRemote()
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
