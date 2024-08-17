package tui

type ViewPosition struct {
	viewName string
	x0       int
	y0       int
	x1       int
	y1       int
	overlaps byte
}

func NewViewPosition(
	viewName string,
	x0 int,
	y0 int,
	x1 int,
	y1 int,
	overlaps byte,
) *ViewPosition {
	return &ViewPosition{
		viewName: viewName,
		x0:       x0,
		y0:       y0,
		x1:       x1,
		y1:       y1,
		overlaps: overlaps,
	}
}

func (v *ViewPosition) SetName(n string) {
	v.viewName = n
}

func (p *ViewPosition) Set() *View {
	return SetView(p.viewName, p.x0, p.y0, p.x1, p.y1, p.overlaps)
}
