package ui

import (
	"errors"
	"log"
	"mynav/pkg/api"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	gui *gocui.Gui
	api *api.Api
	State
}

type State struct {
	help          *HelpState
	confirmation  *ConfirmationDialogState
	toast         *ToastDialogState
	editor        *EditorDialogState
	header        *HeaderState
	workspaceInfo *WorkspaceInfoDialogState
	workspaces    *WorkspacesState
	topics        *TopicsState
	fs            *FsState
	action        *Action
}

func Start() *Action {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	ui := &UI{
		gui: g,
		api: api.NewApi(),
		State: State{
			confirmation:  newConfirmationDialogState(),
			editor:        newEditorDialogState(),
			toast:         newToastDialogState(),
			workspaces:    newWorkspacesState(),
			header:        newHeaderState(),
			workspaceInfo: newWorkspaceInfoDialogState(),
			topics:        newTopicsState(),
			fs:            newFsState(),
		},
	}

	ui.help = ui.newHelpState(ui.getKeyBindings("global"))

	ui.gui.SetManager(gocui.ManagerFunc(ui.renderViews))
	quit := func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}
	ui.keyBinding("").
		setKeybinding("", gocui.KeyCtrlC, quit).
		setKeybinding("", 'q', quit).
		set('?', func() {
			ui.openHelpView(nil)
		})

	err = ui.gui.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}

	return ui.action
}

func (ui *UI) handleUpdate() {
	if ui.api.IsConfigInitialized && !ui.api.IsUpdateAsked() {
		ui.api.SetUpdateAsked()
		update, newTag := ui.api.DetectUpdate()
		if update {
			ui.openConfirmationDialog(func(b bool) {
				if b {
					ui.setActionEnd(ui.api.GetUpdateSystemCmd())
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
		}
	}
}

func (ui *UI) renderViews(g *gocui.Gui) error {
	ui.renderHeaderView()

	ui.handleUpdate()

	ui.renderFsView()
	ui.renderToastDialog()
	ui.renderConfirmationDialog()
	ui.renderWorkspaceInfoDialog()
	ui.renderEditorDialog()
	ui.renderHelpView()

	if ui.action != nil && ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}
