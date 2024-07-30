package ui

type ListRenderer struct {
	selected   int
	startIdx   int
	endIdx     int
	listSize   int
	renderSize int
}

func NewListRenderer(initial int, renderSize int, listSize int) *ListRenderer {
	size := min(listSize, renderSize)
	return &ListRenderer{
		selected:   initial,
		startIdx:   initial,
		endIdx:     initial + size,
		listSize:   listSize,
		renderSize: renderSize,
	}
}

func (lr *ListRenderer) Increment() {
	if lr.selected >= lr.listSize-1 {
		lr.SetSelected(0)
		return
	}

	lr.selected++
	if lr.selected >= lr.endIdx {
		lr.endIdx++
		lr.startIdx++
	}
}

func (lr *ListRenderer) Decrement() {
	if lr.selected <= 0 {
		lr.SetSelected(lr.listSize - 1)
		return
	}
	lr.selected--
	if lr.selected < lr.startIdx {
		lr.endIdx--
		lr.startIdx--
	}
}

func (lr *ListRenderer) SetSelected(idx int) {
	if idx < 0 || idx > lr.listSize {
		return
	}

	size := min(lr.renderSize, lr.listSize)
	lr.selected = idx
	lr.startIdx = min(lr.selected, lr.listSize-size)
	lr.endIdx = lr.startIdx + size
}

func (lr *ListRenderer) ResetSize(newSize int) {
	if newSize != lr.listSize {
		lr.setListSize(newSize)
	}
}

func (lr *ListRenderer) setListSize(listSize int) {
	lr.listSize = listSize
	if listSize == 0 {
		lr.SetSelected(0)
	} else if lr.selected >= listSize {
		lr.SetSelected(listSize - 1)
	} else {
		lr.SetSelected(lr.selected)
	}
}

func (lr *ListRenderer) forEach(f func(idx int)) {
	for i := lr.startIdx; i < lr.endIdx; i++ {
		f(i)
	}
}
