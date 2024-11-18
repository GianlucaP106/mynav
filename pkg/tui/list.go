package tui

type ListRenderer struct {
	selected int
	startIdx int
	endIdx   int
	realSize int
	size     int
}

func NewListRenderer(initial int, renderSize int, listSize int) *ListRenderer {
	size := min(listSize, renderSize)
	return &ListRenderer{
		selected: initial,
		startIdx: initial,
		endIdx:   initial + size,
		realSize: listSize,
		size:     renderSize,
	}
}

func (lr *ListRenderer) Increment() {
	if lr.selected >= lr.realSize-1 {
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
		lr.SetSelected(lr.realSize - 1)
		return
	}
	lr.selected--
	if lr.selected < lr.startIdx {
		lr.endIdx--
		lr.startIdx--
	}
}

func (lr *ListRenderer) SetSelected(idx int) {
	if idx < 0 || idx > lr.realSize {
		return
	}

	size := min(lr.size, lr.realSize)
	lr.selected = idx
	lr.startIdx = min(lr.selected, lr.realSize-size)
	lr.endIdx = lr.startIdx + size
}

func (lr *ListRenderer) ResetSize(newSize int) {
	if newSize != lr.realSize {
		lr.setListSize(newSize)
	}
}

func (lr *ListRenderer) setListSize(listSize int) {
	lr.realSize = listSize
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
