package ui

import (
	"github.com/awesome-gocui/gocui"
)

type EditFunc func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier)

type Editor interface {
	Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier)
}

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

func NewConfirmationEditor(onConfirm func(), onReject func()) gocui.EditorFunc {
	return gocui.EditorFunc(
		func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			switch key {
			case gocui.KeyEnter:
				onConfirm()
			case gocui.KeyEsc:
				onReject()
			}
		})
}

func NewSingleActionEditor(keys []gocui.Key, action func()) gocui.EditorFunc {
	return gocui.EditorFunc(
		func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			for _, k := range keys {
				if key == k {
					action()
					return
				}
			}
		})
}

func NewListRendererEditor(up func(), down func(), enter func(), exit func()) gocui.EditorFunc {
	return gocui.EditorFunc(
		func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			switch {
			case key == gocui.KeyEnter:
				enter()
			case key == gocui.KeyEsc:
				exit()
			case ch == 'j':
				down()
			case ch == 'k':
				up()
			}
		})
}
