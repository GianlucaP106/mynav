package core

import (
	"log"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"os"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/shirou/gopsutil/process"
)

type TmuxController struct {
	portController    *system.PortController
	processController *system.ProcessController
	tmux              *gotmux.Tmux
}

func NewTmuxController(pc *system.PortController, pcc *system.ProcessController) *TmuxController {
	t, err := gotmux.DefaultTmux()
	if err != nil {
		log.Panicln(err)
	}

	tmc := &TmuxController{
		portController:    pc,
		tmux:              t,
		processController: pcc,
	}

	return tmc
}

func (tc *TmuxController) RenameTmuxSession(s *gotmux.Session, newName string) error {
	err := s.Rename(newName)
	if err != nil {
		return err
	}

	events.Emit(events.TmuxSessionChangeEvent)
	return nil
}

func (tc *TmuxController) CreateAndAttachTmuxSession(session string, path string) error {
	s, err := tc.tmux.NewSession(&gotmux.SessionOptions{
		Name:           session,
		StartDirectory: path,
	})
	if err != nil {
		return err
	}

	err = s.Attach()
	if err != nil {
		return err
	}

	events.Emit(events.PortChangeEvent)
	events.Emit(events.TmuxSessionChangeEvent)
	return nil
}

func (tc *TmuxController) AttachTmuxSession(s *gotmux.Session) error {
	err := s.Attach()
	if err != nil {
		return err
	}

	events.Emit(events.PortChangeEvent)
	events.Emit(events.TmuxSessionChangeEvent)
	return nil
}

func (tc *TmuxController) GetTmuxSessionCount() int {
	s, _ := tc.tmux.ListSessions()
	return len(s)
}

func (tc *TmuxController) GetTmuxSessions() []*gotmux.Session {
	sessions, err := tc.tmux.ListSessions()
	if err != nil {
		return []*gotmux.Session{}
	}

	return sessions
}

func (tc *TmuxController) GetTmuxSessionByName(name string) *gotmux.Session {
	session, _ := tc.tmux.GetSessionByName(name)
	return session
}

func (tc *TmuxController) DeleteTmuxSession(s *gotmux.Session) error {
	err := s.Kill()
	if err != nil {
		return err
	}

	events.Emit(events.TmuxSessionChangeEvent)
	events.Emit(events.PortChangeEvent)
	return nil
}

func (tc *TmuxController) KillTmuxServer() error {
	err := tc.tmux.KillServer()
	if err != nil {
		return err
	}

	events.Emit(events.TmuxSessionChangeEvent)
	return nil
}

func (tc *TmuxController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0

	sessions, err := tc.tmux.ListSessions()
	if err != nil {
		return
	}

	for _, s := range sessions {
		sessionCount++
		windowCount += s.Windows
	}

	return
}

func (tc *TmuxController) KillTmuxWindow(w *gotmux.Window) error {
	err := w.Kill()
	if err != nil {
		return err
	}

	events.Emit(events.TmuxSessionChangeEvent)
	return nil
}

func (tc *TmuxController) KillTmuxPane(pane *gotmux.Pane) error {
	err := pane.Kill()
	if err != nil {
		return err
	}

	events.Emit(events.TmuxSessionChangeEvent)
	return nil
}

func (t *TmuxController) GetTmuxSessionChildProcesses(session *gotmux.Session) []*process.Process {
	windows, err := session.ListWindows()
	out := make([]*process.Process, 0)
	if err != nil {
		return out
	}

	for _, window := range windows {
		panes, err := window.ListPanes()
		if err != nil {
			continue
		}

		for _, p := range panes {
			proc, err := process.NewProcess(p.Pid)
			if err != nil {
				continue
			}

			children, err := proc.Children()
			if err != nil {
				continue
			}

			out = append(out, proc)
			out = append(out, children...)
		}
	}

	return out
}

func (tc *TmuxController) GetTmuxSessionByPort(port *system.Port) *gotmux.Session {
	return tc.GetTmuxSessionByChildPid(port.Pid)
}

func (t *TmuxController) GetTmuxSessionByChildPid(pid int) *gotmux.Session {
	for _, session := range t.GetTmuxSessions() {
		windows, err := session.ListWindows()
		if err != nil {
			continue
		}

		for _, w := range windows {
			panes, err := w.ListPanes()
			if err != nil {
				continue
			}

			for _, p := range panes {
				if t.processController.IsProcessChildOf(pid, int(p.Pid)) {
					return session
				}
			}
		}

	}

	return nil
}

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}
