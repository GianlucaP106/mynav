package tui

import (
	"github.com/awesome-gocui/gocui"
)

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

			case key == gocui.KeyCtrlA || key == gocui.KeyCtrlE:
				_, cury := v.Cursor()
				line, err := v.Line(cury)
				if err != nil {
					return
				}

				switch key {
				case gocui.KeyCtrlA: // <-
					v.SetCursor(0, cury)
				case gocui.KeyCtrlE: // ->
					v.SetCursor(len(line), cury)

				}
			case mod == gocui.ModAlt && (ch == 'b' || ch == 'f'):
				curx, cury := v.Cursor()
				curLine, err := v.Line(cury)
				if err != nil {
					return
				}

				indicies := getSkipIndicies(curLine)

			outer:
				for thisIdx, idx := range indicies[:len(indicies)-1] {
					before := idx
					after := indicies[thisIdx+1]

					switch ch {
					case 'b': // <-
						if before < curx && after >= curx {
							v.SetCursor(before, cury)
							break outer
						}
					case 'f': // ->
						if before <= curx && after > curx {
							v.SetCursor(after, cury)
							break outer
						}
					}
				}
			}
		})
}

func getSkipIndicies(line string) []int {
	indicies := []int{}
	spaceBlock := false
	for idx, c := range line {
		if c == ' ' {
			if !spaceBlock {
				spaceBlock = true
				indicies = append(indicies, idx)
			}
		} else {
			spaceBlock = false
		}
	}

	indicies = append([]int{0}, indicies...)
	indicies = append(indicies, len(line))
	return indicies
}
