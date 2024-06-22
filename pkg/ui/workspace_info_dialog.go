package ui

import (
	"fmt"
	"mynav/pkg/core"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const WorkspaceInfoDialogStateName = "WorkspaceInfoDialog"

type WorkspaceInfoDialog struct {
	editor    Editor
	workspace *core.Workspace
	title     string
}

var _ Dialog = &WorkspaceInfoDialog{}

func newWorkspaceInfoDialogState() *WorkspaceInfoDialog {
	return &WorkspaceInfoDialog{}
}

func (w *WorkspaceInfoDialog) Name() string {
	return WorkspaceInfoDialogStateName
}

func (w *WorkspaceInfoDialog) Init(height int) *gocui.View {
	view := SetCenteredView(w.Name(), 100, height, 0)
	view.Title = withSurroundingSpaces(w.title)
	view.TitleColor = gocui.ColorBlue
	view.Editor = w.editor
	view.Editable = true
	FocusView(w.Name())
	return view
}

func (wd *WorkspaceInfoDialog) Open(w *core.Workspace, exit func()) {
	prevView := GetFocusedView()
	wd.editor = NewConfirmationEditor(func() {
		wd.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		exit()
	}, func() {
		wd.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		exit()
	})

	wd.workspace = w
	wd.title = w.Name

	content := wd.formatWorkspaceInfo(wd.workspace)
	wd.Init(len(content))
}

func (wd *WorkspaceInfoDialog) Close() {
	DeleteView(wd.Name())
}

func (wd *WorkspaceInfoDialog) formatWorkspaceInfo(w *core.Workspace) []string {
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
		desc := splitStringByLength(w.Metadata.Description, sizeX)
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

	if s := Api().Tmux.GetTmuxSessionByName(w.Path); s != nil {
		out = append(out, formatItem("Tmux session: ", s.Name)...)
		out = append(out, withSpacePadding(strconv.Itoa(s.NumWindows)+" window(s)", sizeX))
		out = append(out, blankLine(sizeX))
	}

	if pw := Api().Core.GetPortsByWorkspace(w); pw != nil && pw.Len() > 0 {
		ports := ""
		for _, p := range pw {
			ports += p.GetPortStr() + ", "
		}

		ports = trimEnd(ports, 2)
		out = append(out, formatItem("Open Ports: ", ports)...)
		out = append(out, blankLine(sizeX))

	}

	out = append(out, formatItem("Last modified: ", w.GetLastModifiedTimeFormatted())...)
	out = append(out, blankLine(sizeX))

	out = append(out, blankLine(sizeX))
	out = append(out, blankLine(sizeX))

	return out
}

func (w *WorkspaceInfoDialog) Render(ui *UI) error {
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

func splitStringByLength(input string, chunkSize int) []string {
	var chunks []string
	for len(input) > 0 {
		if len(input) >= chunkSize {
			chunks = append(chunks, input[:chunkSize])
			input = input[chunkSize:]
		} else {
			chunks = append(chunks, input)
			break
		}
	}
	return chunks
}
