package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/core"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type WorkspaceInfoDialog struct {
	view      *View
	workspace *core.Workspace
}

func OpenWorkspaceInfoDialog(w *core.Workspace, exit func()) *WorkspaceInfoDialog {
	wd := &WorkspaceInfoDialog{}
	wd.workspace = w
	content := wd.getWorkspaceInfoContent(wd.workspace)
	wd.view = SetCenteredView(constants.WorkspaceInfoDialogName, 100, len(content), 0)

	wd.view.Title = withSurroundingSpaces(wd.workspace.Name)
	wd.view.TitleColor = gocui.ColorBlue
	wd.view.Editable = true

	prevView := GetFocusedView()
	wd.view.Editor = NewConfirmationEditor(func() {
		wd.Close()
		if prevView != nil {
			SetFocusView(prevView.Name())
		}
		exit()
	}, func() {
		wd.Close()
		if prevView != nil {
			SetFocusView(prevView.Name())
		}
		exit()
	})

	wd.view.Clear()
	for _, line := range content {
		fmt.Fprintln(wd.view, line)
	}

	SetFocusView(wd.view.Name())

	return wd
}

func (wd *WorkspaceInfoDialog) Close() {
	wd.view.Delete()
}

func (wd *WorkspaceInfoDialog) getWorkspaceInfoContent(w *core.Workspace) []string {
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
		out = append(out, withSpacePadding(strconv.Itoa(s.Windows)+" window(s)", sizeX))
		out = append(out, blankLine(sizeX))
	}

	out = append(out, formatItem("Last modified: ", w.GetLastModifiedTimeFormatted())...)
	out = append(out, blankLine(sizeX))

	out = append(out, blankLine(sizeX))
	out = append(out, blankLine(sizeX))

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
