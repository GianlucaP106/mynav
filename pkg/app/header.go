package app

import (
	"fmt"
	"mynav/pkg/tui"
	"strconv"
	"sync/atomic"
)

type Header struct {
	// stats views
	lwv *tui.View
	scv *tui.View
	wcv *tui.View
	tcv *tui.View

	// stats data (atomic as they may be modfied by worker)
	lastWorkspace  atomic.Value
	sessionCount   atomic.Int32
	workspaceCount atomic.Int32
	topicCount     atomic.Int32
}

func newHeader() *Header {
	h := &Header{}
	h.lastWorkspace.Store("")
	return h
}

// Refreshes the data.
func (hv *Header) refresh() {
	hv.workspaceCount.Store(int32(a.api.WorkspacesCount()))
	hv.topicCount.Store(int32(a.api.TopicCount()))
	w := a.api.SelectedWorkspace()
	if w != nil {
		hv.lastWorkspace.Store(w.ShortPath())
	} else {
		hv.lastWorkspace.Store("")
	}

	// session count last because it is slow
	hv.sessionCount.Store(int32(a.api.SessionCount()))
}

func (hv *Header) init() {
	// last workspace header config
	hv.lwv = a.ui.SetView(getViewPosition(HeaderView))
	hv.lwv.Title = " Last Workspace "
	a.styleView(hv.lwv)

	// topic count count header config
	hv.tcv = a.ui.SetView(getViewPosition(Header2View))
	hv.tcv.Title = " Topics "
	a.styleView(hv.tcv)

	// workspace count header config
	hv.wcv = a.ui.SetView(getViewPosition(Header3View))
	hv.wcv.Title = " Workspaces "
	a.styleView(hv.wcv)

	// session count header config
	hv.scv = a.ui.SetView(getViewPosition(Header4View))
	hv.scv.Title = " Sessions "
	a.styleView(hv.scv)
}

func (hv *Header) render() {
	hv.renderLastWorkspace()
	hv.renderTopicCount()
	hv.renderWorkspaceCount()
	hv.renderSessionCount()
}

func (hv *Header) renderLastWorkspace() {
	// last workspace header section
	hv.lwv.Clear()
	hv.lwv = a.ui.SetView(getViewPosition(hv.lwv.Name()))
	lastWorkspace := hv.lastWorkspace.Load().(string)
	if lastWorkspace == "" {
		return
	}

	line := lastWorkspace
	line = topicNameColor.Sprint(line)
	fmt.Fprintln(hv.lwv, line)
}

func (hv *Header) renderTopicCount() {
	// topic count count header section
	hv.tcv.Clear()
	hv.tcv = a.ui.SetView(getViewPosition(hv.tcv.Name()))

	count := " " + strconv.Itoa(int(hv.topicCount.Load()))
	count = workspaceNameColor.Sprint(count)
	fmt.Fprintln(hv.tcv, count)
}

func (hv *Header) renderWorkspaceCount() {
	// workspace count count header section
	hv.wcv.Clear()
	hv.wcv = a.ui.SetView(getViewPosition(hv.wcv.Name()))

	count := " " + strconv.Itoa(int(hv.workspaceCount.Load()))
	style := alternateSessionMarkerColor
	count = style.Sprint(count)
	fmt.Fprintln(hv.wcv, count)
}

func (hv *Header) renderSessionCount() {
	// session count header section
	hv.scv.Clear()
	hv.scv = a.ui.SetView(getViewPosition(hv.scv.Name()))

	count := " " + strconv.Itoa(int(hv.sessionCount.Load()))
	s := sessionMarkerColor
	count = s.Sprint(count)
	fmt.Fprintln(hv.scv, count)
}
