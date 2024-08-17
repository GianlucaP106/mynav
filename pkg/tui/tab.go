package tui

import "github.com/awesome-gocui/gocui"

type TabGroup struct {
	focusView func(string)
	Tabs      []*Tab
	Selected  int
}

type Tab struct {
	focusView   func(string)
	Frame       *View
	DefaultView string
	LastView    string
	Views       []*ViewSlot
}

type ViewSlot struct {
	View             *View
	ViewSlotPosition ViewSlotPosition
}

type ViewSlotPosition = uint

const (
	NoPosition ViewSlotPosition = iota
	TopLeftPosition
	TopRightPosition
	BottomLeftPosition
	BottomRightPosition
)

func NewTabGroup(focus func(string), tabs ...*Tab) *TabGroup {
	tg := &TabGroup{
		Tabs:      make([]*Tab, 0),
		focusView: focus,
	}
	tg.Tabs = append(tg.Tabs, tabs...)
	tg.Selected = 0
	return tg
}

func (tg *TabGroup) AddTab(t *Tab) {
	tg.Tabs = append(tg.Tabs, t)
}

func (tg *TabGroup) GetTab(name string) *Tab {
	for _, t := range tg.Tabs {
		if t.Frame.Name() == name {
			return t
		}
	}

	return nil
}

func (tg *TabGroup) IncrementSelectedTab(callback func(*Tab)) {
	if tg.Selected >= len(tg.Tabs)-1 {
		tg.Selected = 0
		tg.FocusTabByIndex(0)
		callback(tg.GetSelectedTab())
	} else {
		tg.FocusTabByIndex(tg.Selected + 1)
		callback(tg.GetSelectedTab())
	}
}

func (tg *TabGroup) DecrementSelectedTab(callback func(*Tab)) {
	if tg.Selected <= 0 {
		newIdx := len(tg.Tabs) - 1
		tg.Selected = newIdx
		tg.FocusTabByIndex(newIdx)
		callback(tg.GetSelectedTab())
	} else {
		tg.FocusTabByIndex(tg.Selected - 1)
		callback(tg.GetSelectedTab())
	}
}

func (tg *TabGroup) GetSelectedTab() *Tab {
	if len(tg.Tabs) > 0 {
		return tg.Tabs[tg.Selected]
	}

	return nil
}

func (tg *TabGroup) FocusTabByIndex(idx int) {
	tg.Selected = idx
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

func (tg *TabGroup) NewTab(name string, defaultView string) *Tab {
	t := &Tab{
		Views:       make([]*ViewSlot, 0),
		DefaultView: defaultView,
		LastView:    defaultView,
		focusView:   tg.focusView,
	}
	x, y := ScreenSize()
	t.Frame = SetCenteredView(name, x, y, 0)
	t.Frame.Frame = false
	tg.Tabs = append(tg.Tabs, t)
	return t
}

func (t *Tab) GetTabView(view string) *View {
	for _, v := range t.Views {
		if v.View.Name() == view {
			return v.View
		}
	}

	return nil
}

func (t *Tab) AddView(v *View, position ViewSlotPosition) {
	t.Views = append(t.Views, NewViewSlot(v, position))
}

func (t *Tab) SendToFront() {
	t.Frame.SendToFront()
	for _, v := range t.Views {
		v.View.SendToFront()
	}

	t.focusView(t.LastView)
}

func (t *Tab) SendToBack() {
	lastView := GetFocusedView()
	if lastView != nil && t.GetTabView(lastView.Name()) != nil {
		t.LastView = lastView.Name()
	} else {
		t.LastView = t.DefaultView
	}

	t.Frame.SendToBack()
	for _, v := range t.Views {
		v.View.SendToBack()
	}
}

type SlotPositionRelation struct {
	keys           []gocui.Key
	targetPosition ViewSlotPosition
}

func (t *Tab) GenerateNavigationKeyBindings() {
	views := map[ViewSlotPosition]*ViewSlot{}
	findView := func(pos ViewSlotPosition) *ViewSlot {
		for _, vs := range t.Views {
			if vs.ViewSlotPosition == pos {
				return vs
			}
		}

		return nil
	}
	views[TopLeftPosition] = findView(TopLeftPosition)
	views[BottomLeftPosition] = findView(BottomLeftPosition)
	views[TopRightPosition] = findView(TopRightPosition)
	views[BottomRightPosition] = findView(BottomRightPosition)

	relationMap := make(map[ViewSlotPosition][]*SlotPositionRelation)
	relationMap[TopLeftPosition] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowDown, gocui.KeyCtrlJ},
			targetPosition: BottomLeftPosition,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowRight, gocui.KeyCtrlL},
			targetPosition: TopRightPosition,
		},
	}

	relationMap[BottomLeftPosition] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowUp, gocui.KeyCtrlK},
			targetPosition: TopLeftPosition,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowRight, gocui.KeyCtrlL},
			targetPosition: BottomRightPosition,
		},
	}

	relationMap[BottomRightPosition] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowUp, gocui.KeyCtrlK},
			targetPosition: TopRightPosition,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowLeft, gocui.KeyCtrlH},
			targetPosition: BottomLeftPosition,
		},
	}

	relationMap[TopRightPosition] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowDown, gocui.KeyCtrlJ},
			targetPosition: BottomRightPosition,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowLeft, gocui.KeyCtrlH},
			targetPosition: TopLeftPosition,
		},
	}

	for pos, slotRelation := range relationMap {
		viewSlot := views[pos]
		if viewSlot == nil {
			continue
		}

		kbb := NewKeybindingBuilder(viewSlot.View)
		for _, relation := range slotRelation {
			targetSlot := views[relation.targetPosition]
			if targetSlot == nil {
				continue
			}

			viewName := targetSlot.View.Name()
			for _, key := range relation.keys {
				kbb.Set(key, "Focus "+viewName, func() {
					t.focusView(viewName)
				})
			}
		}
	}
}

func NewViewSlot(v *View, position ViewSlotPosition) *ViewSlot {
	return &ViewSlot{
		View:             v,
		ViewSlotPosition: position,
	}
}
