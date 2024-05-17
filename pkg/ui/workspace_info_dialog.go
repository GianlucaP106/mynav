package ui

import (
	"fmt"
	"mynav/pkg/api"
	"mynav/pkg/utils"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type WorkspaceInfoDialogState struct {
	editor    Editor
	workspace *api.Workspace
	viewName  string
	title     string
	active    bool
}

func newWorkspaceInfoDialogState() *WorkspaceInfoDialogState {
	return &WorkspaceInfoDialogState{
		viewName: "WorkspaceInfoDialog",
	}
}

func (ui *UI) initWorkspaceInfoDialog(height int) *gocui.View {
	view := ui.setCenteredView(ui.workspaceInfo.viewName, 100, height, 0)
	view.Title = withSurroundingSpaces(ui.workspaceInfo.title)
	view.TitleColor = gocui.ColorBlue
	view.Editor = ui.workspaceInfo.editor
	view.Editable = true
	return view
}

func (ui *UI) openWorkspaceInfoDialog(w *api.Workspace) {
	ui.workspaceInfo.editor = newConfirmationEditor(func() {
		ui.closeWorkspaceInfoDialog()
	}, func() {
		ui.closeWorkspaceInfoDialog()
	})
	ui.workspaceInfo.workspace = w
	ui.workspaceInfo.title = w.Name
	ui.workspaceInfo.active = true
}

func (ui *UI) closeWorkspaceInfoDialog() {
	ui.workspaceInfo.active = false
	ui.gui.DeleteView(ui.workspaceInfo.viewName)
}

func (ui *UI) formatWorkspaceInfo(w *api.Workspace) []string {
	sizeX := 100

	formatItem := func(title string, content string) []string {
		return []string{
			withSpacePadding(color.Blue.Sprint(title), sizeX),
			withSpacePadding(content, sizeX),
		}
	}

	description := func() []string {
		out := []string{}
		if w.Metadata.Description == "" {
			return out
		}
		out = append(out, withSpacePadding(color.Blue.Sprint("Description: "), sizeX))
		desc := utils.SplitStringByLength(w.Metadata.Description, sizeX)
		out = append(out, desc...)
		return out
	}()

	out := []string{}
	out = append(out, blankLine(sizeX))

	for _, line := range description {
		out = append(out, color.White.Sprint(line))
	}
	out = append(out, blankLine(sizeX))

	// TODO: handle error
	remote, _ := w.GetGitRemote()
	if remote != "" {
		out = append(out, formatItem("Git remote: ", remote)...)
		out = append(out, blankLine(sizeX))
	}

	if w.Metadata.TmuxSession != nil && w.Metadata.TmuxSession.Name != "" {
		out = append(out, formatItem("Tmux session: ", w.Metadata.TmuxSession.Name)...)
		out = append(out, withSpacePadding(strconv.Itoa(w.Metadata.TmuxSession.NumWindows)+" window(s)", sizeX))
		out = append(out, blankLine(sizeX))
	}

	out = append(out, formatItem("Last modified: ", w.GetLastModifiedTimeFormatted())...)
	out = append(out, blankLine(sizeX))

	out = append(out, blankLine(sizeX))
	out = append(out, blankLine(sizeX))

	return out
}

func (ui *UI) renderWorkspaceInfoDialog() {
	if !ui.workspaceInfo.active {
		return
	}

	content := ui.formatWorkspaceInfo(ui.workspaceInfo.workspace)
	view := ui.initWorkspaceInfoDialog(len(content))
	ui.focusView(ui.workspaceInfo.viewName)

	view.Clear()
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
}
