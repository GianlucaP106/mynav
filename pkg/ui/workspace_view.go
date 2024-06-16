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

const WorkspacesViewName = "WorkspacesView"

type WorkspacesView struct {
	listRenderer *ListRenderer
	tv           *TopicsView
	search       string
	workspaces   api.Workspaces
}

var _ View = &WorkspacesView{}

func newWorkspacesView(tv *TopicsView) *WorkspacesView {
	workspace := &WorkspacesView{
		tv: tv,
	}
	return workspace
}

func (wv *WorkspacesView) selectWorkspaceByShortPath(shortPath string) {
	for idx, w := range wv.workspaces {
		if w.ShortPath() == shortPath {
			wv.listRenderer.setSelected(idx)
		}
	}
}

func (wv *WorkspacesView) refreshWorkspaces() {
	tv := wv.tv
	var workspaces api.Workspaces
	if selectedTopic := tv.getSelectedTopic(); selectedTopic != nil {
		workspaces = Api().GetWorkspaces().ByTopic(selectedTopic)
	} else {
		workspaces = make(api.Workspaces, 0)
	}

	if wv.search != "" {
		workspaces = workspaces.FilterByNameContaining(wv.search)
	}

	wv.workspaces = workspaces

	if wv.listRenderer != nil {
		newListSize := len(wv.workspaces)
		if wv.listRenderer.listSize != newListSize {
			wv.listRenderer.setListSize(newListSize)
		}
	}
}

func (wv *WorkspacesView) getSelectedWorkspace() *api.Workspace {
	idx := wv.listRenderer.selected
	if idx >= len(wv.workspaces) || idx < 0 {
		return nil
	}
	return wv.workspaces[idx]
}

func (wv *WorkspacesView) RequiresManager() bool {
	return false
}

func (wv *WorkspacesView) Name() string {
	return WorkspacesViewName
}

func (wv *WorkspacesView) Init(ui *UI) {
	if GetInternalView(wv.Name()) != nil {
		return
	}

	view := SetViewLayout(wv.Name())

	view.Title = withSurroundingSpaces("Workspaces")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorBlue

	// TODO: change
	if Api().GetSelectedWorkspace() != nil {
		ui.FocusWorkspacesView()
	} else {
		ui.FocusTopicsView()
	}

	_, sizeY := view.Size()
	wv.listRenderer = newListRenderer(0, sizeY, 0)
	wv.refreshWorkspaces()

	if selectedWorkspace := Api().GetSelectedWorkspace(); selectedWorkspace != nil {
		wv.selectWorkspaceByShortPath(selectedWorkspace.ShortPath())
	}

	KeyBinding(wv.Name()).
		set('j', func() {
			wv.listRenderer.increment()
		}).
		set('k', func() {
			wv.listRenderer.decrement()
		}).
		set(gocui.KeyArrowDown, func() {
			ui.FocusTmuxView()
		}).
		set(gocui.KeyArrowLeft, func() {
			ui.FocusTopicsView()
		}).
		set(gocui.KeyEsc, func() {
			if wv.search != "" {
				wv.search = ""
				wv.refreshWorkspaces()
				return
			}

			ui.FocusTopicsView()
		}).
		set('s', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}
			GetDialog[*WorkspaceInfoDialogState](ui).Open(curWorkspace, func() {
				ui.FocusWorkspacesView()
			})
		}).
		set('g', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().CloneRepo(s, curWorkspace); err != nil {
					GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
						ui.FocusWorkspacesView()
					})
				}
				ui.FocusWorkspacesView()
			}, func() {
				ui.FocusWorkspacesView()
			}, "Git repo URL", Small)
		}).
		set('/', func() {
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				wv.search = s
				wv.refreshWorkspaces()
				ui.FocusWorkspacesView()
			}, func() {
				ui.FocusWorkspacesView()
			}, "Search", Small)
		}).
		setKeybinding(wv.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			if utils.IsTmuxSession() {
				ui.setAction(Api().GetWorkspaceNvimCmd(curWorkspace))
				return gocui.ErrQuit
			}

			command := Api().GetCreateOrAttachTmuxSessionCmd(curWorkspace)
			ui.setAction(command)

			return gocui.ErrQuit
		}).
		setKeybinding(wv.Name(), 'v', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			Api().SetSelectedWorkspace(curWorkspace)
			ui.setAction(utils.NvimCmd(curWorkspace.Path))
			return gocui.ErrQuit
		}).
		setKeybinding(wv.Name(), 'm', func(g *gocui.Gui, v *gocui.View) error {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return nil
			}

			openTermCmd, err := utils.GetOpenTerminalCmd(curWorkspace.Path)
			if err != nil {
				GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
					ui.FocusWorkspacesView()
				})
				return nil
			}

			ui.setAction(openTermCmd)
			return gocui.ErrQuit
		}).
		set('D', func() {
			if Api().GetWorkspacesByTopicCount(wv.tv.getSelectedTopic()) <= 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					curWorkspace := wv.getSelectedWorkspace()
					Api().DeleteWorkspace(curWorkspace)

					// HACK: same as below
					wv.tv.listRenderer.setSelected(0)
					ui.RefreshMainView()
				}
				ui.FocusWorkspacesView()
			}, "Are you sure you want to delete this workspace?")
		}).
		set('r', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().RenameWorkspace(curWorkspace, s); err != nil {
					GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
						ui.FocusWorkspacesView()
					})
					return
				}
				ui.FocusWorkspacesView()
			}, func() {
				ui.FocusWorkspacesView()
			}, "New workspace name", Small)
		}).
		set('e', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(desc string) {
				if desc != "" {
					Api().SetDescription(desc, curWorkspace)
				}
				ui.FocusWorkspacesView()
			}, func() {
				ui.FocusWorkspacesView()
			}, "Description", Large)
		}).
		set('a', func() {
			tv := wv.tv
			curTopic := tv.getSelectedTopic()
			GetDialog[*EditorDialog](ui).Open(func(name string) {
				if _, err := Api().CreateWorkspace(name, curTopic); err != nil {
					GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
						ui.FocusWorkspacesView()
					})
					return
				}

				// HACK: when there a is a new workspace
				// This will result in the workspace and the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.listRenderer.setSelected(0)
				wv.listRenderer.setSelected(0)
				ui.RefreshMainView()
				ui.FocusWorkspacesView()
			}, func() {
				ui.FocusWorkspacesView()
			}, "Workspace name ", Small)
		}).
		set('x', func() {
			curWorkspace := wv.getSelectedWorkspace()
			if curWorkspace == nil {
				return
			}

			if Api().GetTmuxSessionByWorkspace(curWorkspace) != nil {
				GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
					if b {
						Api().DeleteWorkspaceTmuxSession(curWorkspace)
						ui.RefreshMainView()
					}
					ui.FocusWorkspacesView()
				}, "Are you sure you want to delete the tmux session?")
			}
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(workspaceKeyBindings, func() {
				ui.FocusWorkspacesView()
			})
		})
}

func (wv *WorkspacesView) formatWorkspaceRow(workspace *api.Workspace, selected bool) []string {
	sizeX, _ := GetInternalView(wv.Name()).Size()
	style, _ := func() (color.Style, string) {
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

	description := ""
	if workspace.GetDescription() != "" {
		description = withSpacePadding("Description: "+workspace.GetDescription(), fifth)
	}

	url := withSpacePadding(gitRemote, fifth)
	time := withSpacePadding(lastModTime, fifth)

	name := withSurroundingSpaces(workspace.Name)
	tmux := func() string {
		if tm := Api().GetTmuxSessionByWorkspace(workspace); tm != nil {
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

	if description == "" {
		name = withSpacePadding(name, 2*fifth)
	} else {
		name = withSpacePadding(name, fifth)
	}

	line := style.Sprint(name+description+url+time) + tmux + style.Sprint(strings.Repeat(" ", fifth+5)) // +5 for extra padding
	return []string{
		line,
	}
}

func (wv *WorkspacesView) Render(ui *UI) error {
	view := GetInternalView(wv.Name())
	if view == nil {
		wv.Init(ui)
		view = GetInternalView(wv.Name())
	}

	if wv.search != "" {
		view.Subtitle = withSurroundingSpaces("Searching: " + wv.search)
	} else {
		view.Subtitle = ""
	}

	view.Clear()
	content := func() []string {
		if wv.workspaces == nil {
			return []string{}
		}

		out := make([]string, 0)
		wv.listRenderer.forEach(func(i int) {
			w := wv.workspaces[i]
			selected := i == wv.listRenderer.selected && GetFocusedView().Name() == wv.Name()

			// TODO: https://github.com/GianlucaP106/mynav/issues/18
			workspace := wv.formatWorkspaceRow(w, selected)
			out = append(out, workspace...)
		})

		return out
	}()
	for _, line := range content {
		fmt.Fprintln(view, line)
	}

	return nil
}
