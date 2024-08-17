package ui

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/tui"
	"strconv"

	"github.com/gookit/color"
)

type workspaceInfoDialog struct {
	view      *tui.View
	workspace *core.Workspace
}

func openWorkspaceInfoDialog(w *core.Workspace, exit func()) *workspaceInfoDialog {
	wd := &workspaceInfoDialog{}
	wd.workspace = w
	content := wd.getWorkspaceInfoContent(wd.workspace)
	wd.view = tui.SetCenteredView(WorkspaceInfoDialog, 100, len(content), 0)

	wd.view.Title = tui.WithSurroundingSpaces(wd.workspace.Name)
	wd.view.FrameColor = onFrameColor
	styleView(wd.view)
	wd.view.Editable = true

	prevView := tui.GetFocusedView()
	wd.view.Editor = tui.NewConfirmationEditor(func() {
		wd.close()
		if prevView != nil {
			tui.SetFocusView(prevView.Name())
		}
		exit()
	}, func() {
		wd.close()
		if prevView != nil {
			tui.SetFocusView(prevView.Name())
		}
		exit()
	})

	wd.view.Clear()
	for _, line := range content {
		fmt.Fprintln(wd.view, line)
	}

	tui.SetFocusView(wd.view.Name())

	return wd
}

func (wd *workspaceInfoDialog) close() {
	wd.view.Delete()
}

func (wd *workspaceInfoDialog) getWorkspaceInfoContent(w *core.Workspace) []string {
	sizeX := 100

	formatItem := func(title string, content string) []string {
		return []string{
			tui.WithSpaces(color.Blue.Sprint(title), sizeX),
			tui.WithSpaces(content, sizeX),
		}
	}

	description := func() []string {
		out := []string{}
		if w.Metadata.Description == "" {
			return out
		}
		out = append(out, tui.WithSpaces(color.Blue.Sprint("Description: "), sizeX))
		desc := splitStringByLength(w.Metadata.Description, sizeX)
		out = append(out, desc...)
		return out
	}()

	out := []string{}
	out = append(out, tui.BlankLine(sizeX))

	for _, line := range description {
		out = append(out, color.White.Sprint(line))
	}
	out = append(out, tui.BlankLine(sizeX))

	// TODO: handle error
	remote, _ := w.GetGitRemote()
	if remote != "" {
		out = append(out, formatItem("Git remote: ", remote)...)
		out = append(out, tui.BlankLine(sizeX))
	}

	if s := getApi().Tmux.GetTmuxSessionByName(w.Path); s != nil {
		out = append(out, formatItem("Tmux session: ", s.Name)...)
		out = append(out, tui.WithSpaces(strconv.Itoa(s.Windows)+" window(s)", sizeX))
		out = append(out, tui.BlankLine(sizeX))
	}

	out = append(out, formatItem("Last modified: ", w.GetLastModifiedTimeFormatted())...)
	out = append(out, tui.BlankLine(sizeX))

	out = append(out, tui.BlankLine(sizeX))
	out = append(out, tui.BlankLine(sizeX))

	return out
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
