package tui

import "github.com/awesome-gocui/gocui"

func NewSimpleEditor(onEnter func(string), onEsc func(), onType func(string)) gocui.EditorFunc {
	return gocui.EditorFunc(
		func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			switch {
			case ch != 0 && mod == 0:
				v.EditWrite(ch)
				if onType != nil {
					onType(v.Buffer())
				}
			case key == gocui.KeySpace:
				v.EditWrite(' ')
				if onType != nil {
					onType(v.Buffer())
				}
			case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
				v.EditDelete(true)
				if onType != nil {
					onType(v.Buffer())
				}
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
