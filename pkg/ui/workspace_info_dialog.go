package ui

import (
	"fmt"
	"mynav/pkg/api"
	"mynav/pkg/utils"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const WorkspaceInfoDialogStateName = "WorkspaceInfoDialog"

type WorkspaceInfoDialogState struct {
	editor    Editor
	workspace *api.Workspace
	title     string
}

var _ Dialog = &WorkspaceInfoDialogState{}

func newWorkspaceInfoDialogState() *WorkspaceInfoDialogState {
	return &WorkspaceInfoDialogState{}
}

func (w *WorkspaceInfoDialogState) Name() string {
	return WorkspaceInfoDialogStateName
}

func (w *WorkspaceInfoDialogState) Init(height int) *gocui.View {
	view := SetCenteredView(w.Name(), 100, height, 0)
	view.Title = withSurroundingSpaces(w.title)
	view.TitleColor = gocui.ColorBlue
	view.Editor = w.editor
	view.Editable = true
	FocusView(w.Name())
	return view
}

func (wd *WorkspaceInfoDialogState) Open(w *api.Workspace, exit func()) {
	wd.editor = NewConfirmationEditor(func() {
		wd.Close()
		exit()
	}, func() {
		wd.Close()
		exit()
	})

	wd.workspace = w
	wd.title = w.Name

	content := wd.formatWorkspaceInfo(wd.workspace)
	wd.Init(len(content))
}

func (wd *WorkspaceInfoDialogState) Close() {
	DeleteView(wd.Name())
}

func (wd *WorkspaceInfoDialogState) formatWorkspaceInfo(w *api.Workspace) []string {
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

	if s := Api().GetTmuxSessionByWorkspace(w); s != nil {
		out = append(out, formatItem("Tmux session: ", s.Name)...)
		out = append(out, withSpacePadding(strconv.Itoa(s.NumWindows)+" window(s)", sizeX))
		out = append(out, blankLine(sizeX))
	}

	out = append(out, formatItem("Last modified: ", w.GetLastModifiedTimeFormatted())...)
	out = append(out, blankLine(sizeX))

	out = append(out, blankLine(sizeX))
	out = append(out, blankLine(sizeX))

	return out
}

func (w *WorkspaceInfoDialogState) Render(ui *UI) error {
	view := GetInternalView(WorkspaceInfoDialogStateName)
	if view == nil {
		return nil
	}

	content := w.formatWorkspaceInfo(w.workspace)
	view.Clear()
	for _, line := range content {
		fmt.Fprintln(view, line)
	}

	return nil
}
