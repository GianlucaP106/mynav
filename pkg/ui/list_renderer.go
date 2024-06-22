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
		lr.setSelected(0)
		return
	}

	lr.selected++
	if lr.selected >= lr.endIdx {
		lr.endIdx++
		lr.startIdx++
	}
}

func (lr *ListRenderer) decrement() {
	if lr.selected <= 0 {
		lr.setSelected(lr.listSize - 1)
		return
	}
	lr.selected--
	if lr.selected < lr.startIdx {
		lr.endIdx--
		lr.startIdx--
	}
}

func (lr *ListRenderer) setSelected(idx int) {
	if idx < 0 || idx > lr.listSize {
		return
	}

	size := min(lr.renderSize, lr.listSize)
	lr.selected = idx
	lr.startIdx = min(lr.selected, lr.listSize-size)
	lr.endIdx = lr.startIdx + size
}

func (lr *ListRenderer) setListSize(listSize int) {
	lr.listSize = listSize
	if listSize == 0 {
		lr.setSelected(0)
	} else if lr.selected >= listSize {
		lr.setSelected(listSize - 1)
	} else {
		lr.setSelected(lr.selected)
	}
}

func (lr *ListRenderer) forEach(f func(idx int)) {
	for i := lr.startIdx; i < lr.endIdx; i++ {
		f(i)
	}
}

func (lr *ListRenderer) resetSize(newSize int) {
	if newSize != lr.listSize {
		lr.setListSize(newSize)
	}
}
