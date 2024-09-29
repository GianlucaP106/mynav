package ui

import (
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type settingsDialog struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[string]
}
type setting struct {
	set      func(string)
	clear    func()
	isToggle bool
}

func openSettingsDialog() *settingsDialog {
	s := &settingsDialog{}
	screenX, screenY := tui.ScreenSize()
	s.view = tui.SetCenteredView(SettingsDialog, screenX/2, screenY/2, 0)

	s.view.Title = tui.WithSurroundingSpaces("Settings")
	s.view.FrameColor = onFrameColor
	ui.styleView(s.view)

	viewX, viewY := s.view.Size()
	s.tableRenderer = tui.NewTableRenderer[string]()
	s.tableRenderer.InitTable(viewX, viewY, []string{
		"Setting",
		"Value",
	}, []float64{
		0.5,
		0.5,
	})

	settings := map[string]*setting{
		"CustomWorkspaceOpener": {
			set: func(s string) {
				api().GlobalConfiguration.SetCustomWorkspaceOpenerCmd(s)
			},
			clear: func() {
				api().GlobalConfiguration.SetCustomWorkspaceOpenerCmd("")
			},
		},
		"TerminalOpener": {
			set: func(s string) {
				api().GlobalConfiguration.SetTerminalOpenerCmd(s)
			},
			clear: func() {
				api().GlobalConfiguration.SetTerminalOpenerCmd("")
			},
		},
		"EnableGithubTab": {
			set: func(s string) {
				api().GlobalConfiguration.SetGithubTabEnabled(s == "true")
			},
			clear: func() {
				api().GlobalConfiguration.SetGithubTabEnabled(false)
			},
			isToggle: true,
		},
	}

	prevView := tui.GetFocusedView()
	s.view.KeyBinding().
		Set('j', "Move down", func() {
			s.tableRenderer.Down()
			s.render()
		}).
		Set('k', "Move up", func() {
			s.tableRenderer.Up()
			s.render()
		}).
		Set(gocui.KeyEnter, "Change setting", func() {
			cmd := api().GlobalConfiguration.GetCustomWorkspaceOpenerCmd()
			cmdStr := strings.Join(cmd, " ")
			setting := settings[s.getSelectedSetting()]
			if setting == nil {
				return
			}

			if setting.isToggle {
				openConfirmationDialog(func(b bool) {
					if b {
						setting.set("true")
					} else {
						setting.set("false")
					}
					s.refresh()
				}, "Would you like to enable this setting? (Confirm to enable and cancel to disable)")
			} else {
				openEditorDialogWithDefaultValue(func(str string) {
					if str == "" {
						return
					}

					setting.set(str)
					s.refresh()
				}, func() {}, "Change setting", smallEditorSize, cmdStr)
			}
		}).
		Set('D', "Set back to default", func() {
			openConfirmationDialog(func(b bool) {
				if !b {
					return
				}

				setting := settings[s.getSelectedSetting()]
				if setting != nil {
					setting.clear()
				}

				s.refresh()
			}, "Are you sure you want to set this setting back to default?")
		}).
		Set(gocui.KeyEsc, "Close settings", func() {
			s.close()
			if prevView != nil {
				prevView.Focus()
			}
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(s.view.GetKeybindings(), func() {})
		})

	s.view.Focus()
	s.refresh()
	return s
}

func (s *settingsDialog) getSelectedSetting() string {
	_, value := s.tableRenderer.GetSelectedRow()
	if value != nil {
		return *value
	}

	return ""
}

func (s *settingsDialog) refresh() {
	cmd := api().GlobalConfiguration.GetCustomWorkspaceOpenerCmd()
	cmdStr := "tmux/nvim (Default)"
	if len(cmd) > 0 {
		cmdStr = strings.Join(cmd, " ")
	}

	terminalOpener := api().GlobalConfiguration.GetTerminalOpenerCmd()
	defaultOpenTerminalCmd, _ := system.GetOpenTerminalCmd()
	terminalOpenerStr := ""
	if len(terminalOpener) > 0 {
		terminalOpenerStr = strings.Join(terminalOpener, " ")
	} else if len(defaultOpenTerminalCmd) > 0 {
		terminalOpenerStr = strings.Join(defaultOpenTerminalCmd, " ")
		terminalOpenerStr = terminalOpenerStr + " (Default)"
	}

	githubEnabled := "No"
	if api().GlobalConfiguration.GetGithubTabEnabled() {
		githubEnabled = "Yes"
	}

	rows := make([][]string, 0)
	rowValues := make([]string, 0)

	rows = append(rows, []string{
		"Worspace Opener",
		cmdStr,
	})
	rowValues = append(rowValues, "CustomWorkspaceOpener")

	rows = append(rows, []string{
		"Terminal Opener command",
		terminalOpenerStr,
	})
	rowValues = append(rowValues, "TerminalOpener")

	rows = append(rows, []string{
		"Enable Github tab",
		githubEnabled,
	})
	rowValues = append(rowValues, "EnableGithubTab")

	s.tableRenderer.FillTable(rows, rowValues)
	s.render()
}

func (s *settingsDialog) close() {
	s.view.Delete()
}

func (s *settingsDialog) render() {
	s.view.Clear()
	s.tableRenderer.Render(s.view)
}
