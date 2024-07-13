package ui

type TabGroup struct {
	Tabs     []*Tab
	Selected int
}

type Tab struct {
	Frame       *View
	DefaultView string
	LastView    string
	Views       []*View
}

func NewTabGroup(tabs []*Tab) *TabGroup {
	tg := &TabGroup{}
	tg.Tabs = tabs
	tg.Selected = 0
	return tg
}

func (tg *TabGroup) GetTab(name string) *Tab {
	for _, t := range tg.Tabs {
		if t.Frame.Name() == name {
			return t
		}
	}

	return nil
}

func (tg *TabGroup) IncrementSelectedTab() {
	if tg.Selected >= len(tg.Tabs)-1 {
		return
	}

	tg.Selected++
	tg.FocusTabByIndex(tg.Selected)
}

func (tg *TabGroup) DecrementSelectedTab() {
	if tg.Selected <= 0 {
		return
	}

	tg.Selected--
	tg.FocusTabByIndex(tg.Selected)
}

func (tg *TabGroup) GetSelectedTab() *Tab {
	if len(tg.Tabs) > 0 {
		return tg.Tabs[tg.Selected]
	}

	return nil
}

func (tg *TabGroup) FocusTabByIndex(idx int) {
	for i, t := range tg.Tabs {
		if i == idx {
			// to allow for views that are in multiple tabs
			defer t.SendToFront()
		} else {
			t.SendToBack()
		}
	}
}

func (tg *TabGroup) FocusTab(tabName string) {
	idx := -1
	for i, t := range tg.Tabs {
		if t.Frame.Name() == tabName {
			idx = i
		}
	}

	if idx >= 0 {
		tg.FocusTabByIndex(idx)
	}
}

func NewTab(name string, defaultView string) *Tab {
	t := &Tab{
		Views:       make([]*View, 0),
		DefaultView: defaultView,
		LastView:    defaultView,
	}
	x, y := ScreenSize()
	t.Frame = SetCenteredView(name, x, y, 0)
	t.Frame.Frame = false
	return t
}

func (t *Tab) GetTabView(view string) *View {
	for _, v := range t.Views {
		if v.Name() == view {
			return v
		}
	}

	return nil
}

func (t *Tab) AddView(v Viewable) {
	t.Views = append(t.Views, v.View())
}

func (t *Tab) SendToFront() {
	SendViewToFront(t.Frame)
	for _, v := range t.Views {
		SendViewToFront(v)
	}

	FocusView(t.LastView)
}

func (t *Tab) SendToBack() {
	lastView := GetFocusedView()
	if lastView != nil && t.GetTabView(lastView.Name()) != nil {
		t.LastView = lastView.Name()
	} else {
		t.LastView = t.DefaultView
	}

	SendViewToBack(t.Frame)
	for _, v := range t.Views {
		SendViewToBack(v)
	}
}
