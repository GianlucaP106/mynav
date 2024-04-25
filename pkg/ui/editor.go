package ui

import (
	"github.com/awesome-gocui/gocui"
)

type EditFunc func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier)

type Editor interface {
	Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier)
}

func simpleEditorFactory(onEnter func(string), onEsc func()) EditFunc {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
		}
	}
}

func confirmationEditorFactory(onEnter func(), onEsc func()) EditFunc {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch key {
		case gocui.KeyEnter:
			onEnter()
		case gocui.KeyEsc:
			onEsc()
		}
	}
}
