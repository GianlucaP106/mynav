package ui

type ListRenderer struct {
	selected   int
	startIdx   int
	endIdx     int
	listSize   int
	renderSize int
}

func newListRenderer(initial int, renderSize int, listSize int) *ListRenderer {
	size := min(listSize, renderSize)
	return &ListRenderer{
		selected:   initial,
		startIdx:   initial,
		endIdx:     initial + size,
		listSize:   listSize,
		renderSize: renderSize,
	}
}

func (lr *ListRenderer) increment() {
	if lr.selected >= lr.listSize-1 {
		return
	}
	lr.setSelected(lr.selected + 1)
	if lr.selected >= lr.endIdx {
		lr.endIdx++
		lr.startIdx++
	}
}

func (lr *ListRenderer) decrement() {
	if lr.selected <= 0 {
		return
	}
	lr.setSelected(lr.selected - 1)
	if lr.selected < lr.startIdx {
		lr.endIdx--
		lr.startIdx--
	}
}

func (lr *ListRenderer) setSelected(idx int) {
	lr.selected = idx
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (lr *ListRenderer) setListSize(listSize int) {
	newEnd := min(listSize, lr.renderSize)
	if lr.selected >= newEnd {
		lr.selected = newEnd - 1
	}
	if lr.selected < 0 {
		lr.setSelected(0)
	}
	lr.endIdx = newEnd
	lr.listSize = listSize
}

func (lr *ListRenderer) forEach(f func(idx int)) {
	for i := lr.startIdx; i < lr.endIdx; i++ {
		f(i)
	}
}
