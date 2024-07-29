package ui

import "github.com/awesome-gocui/gocui"

type TabGroup struct {
	Tabs     []*Tab
	Selected int
}

type Tab struct {
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
	None ViewSlotPosition = iota
	TopLeft
	TopRight
	BottomLeft
	BottomRight
)

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

func NewTab(name string, defaultView string) *Tab {
	t := &Tab{
		Views:       make([]*ViewSlot, 0),
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
		if v.View.Name() == view {
			return v.View
		}
	}

	return nil
}

func (t *Tab) AddView(v Viewable, position ViewSlotPosition) {
	t.Views = append(t.Views, NewViewSlot(v.View(), position))
}

func (t *Tab) SendToFront() {
	t.Frame.SendToFront()
	for _, v := range t.Views {
		v.View.SendToFront()
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
	views[TopLeft] = findView(TopLeft)
	views[BottomLeft] = findView(BottomLeft)
	views[TopRight] = findView(TopRight)
	views[BottomRight] = findView(BottomRight)

	relationMap := make(map[ViewSlotPosition][]*SlotPositionRelation)
	relationMap[TopLeft] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowDown, gocui.KeyCtrlJ},
			targetPosition: BottomLeft,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowRight, gocui.KeyCtrlL},
			targetPosition: TopRight,
		},
	}

	relationMap[BottomLeft] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowUp, gocui.KeyCtrlK},
			targetPosition: TopLeft,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowRight, gocui.KeyCtrlL},
			targetPosition: BottomRight,
		},
	}

	relationMap[BottomRight] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowUp, gocui.KeyCtrlK},
			targetPosition: TopRight,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowLeft, gocui.KeyCtrlH},
			targetPosition: BottomLeft,
		},
	}

	relationMap[TopRight] = []*SlotPositionRelation{
		{
			keys:           []gocui.Key{gocui.KeyArrowDown, gocui.KeyCtrlJ},
			targetPosition: BottomRight,
		},
		{
			keys:           []gocui.Key{gocui.KeyArrowLeft, gocui.KeyCtrlH},
			targetPosition: TopLeft,
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
				kbb.set(key, "Focus "+viewName, func() {
					FocusView(viewName)
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
