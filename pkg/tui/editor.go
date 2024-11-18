package tui

import "github.com/awesome-gocui/gocui"

func NewSimpleEditor(onEnter func(string), onEsc func()) gocui.EditorFunc {
	return gocui.EditorFunc(
		func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			switch {
			case ch != 0 && mod == 0:
				v.EditWrite(ch)
			case key == gocui.KeySpace:
				v.EditWrite(' ')
			case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
				v.EditDelete(true)
			case key == gocui.KeyEsc:
				onEsc()
			case key == gocui.KeyEnter:
				onEnter(v.Buffer())
			case key == gocui.KeyArrowLeft:
				v.MoveCursor(-1, 0)
			case key == gocui.KeyArrowRight:
				v.MoveCursor(1, 0)
			}
		})
}
