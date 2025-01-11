package tui

import (
	"fmt"
	"io"
	"log"
	"math"
	"sync"

	"github.com/gookit/color"
)

type (
	TableRenderer[T any] struct {
		table        *Table[T]
		listRenderer *ListRenderer
		mu           sync.RWMutex
	}

	Table[T any] struct {
		Title          *TableTitle
		Rows           []*TableRow[T]
		ColProportions []float64

		Width  int
		Height int
	}

	TableTitle struct {
		Titles []string
		Styles []color.Style
	}

	TableRow[T any] struct {
		Value    T
		Cols     []string
		Selected bool
	}
)

func NewTableRenderer[T any]() *TableRenderer[T] {
	return &TableRenderer[T]{}
}

func (tr *TableRenderer[T]) Init(width int, height int, titles []string, colProportions []float64) {
	if len(titles) != len(colProportions) {
		log.Panicln("the number of titles and col proportions should be the same")
	}
	tr.listRenderer = NewListRenderer(0, height-1, 0)

	defaultStyles := []color.Style{}
	for range titles {
		defaultStyles = append(defaultStyles, color.Secondary.Style)
	}

	tr.table = &Table[T]{
		Title: &TableTitle{
			Titles: titles,
			Styles: defaultStyles,
		},
		ColProportions: colProportions,
		Rows:           make([]*TableRow[T], 0),

		Width:  width,
		Height: height,
	}
}

func (tr *TableRenderer[T]) SetStyles(colors []color.Style) {
	tr.table.Title.Styles = colors
}

func (tr *TableRenderer[T]) Size() int {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return len(tr.table.Rows)
}

func (tr *TableRenderer[T]) SelectedRow() (idx int, value *T) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	if len(tr.table.Rows) == 0 {
		return 0, nil
	}
	return tr.listRenderer.selected, &tr.table.Rows[tr.listRenderer.selected].Value
}

func (tr *TableRenderer[T]) SelectRow(idx int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.listRenderer.SetSelected(idx)
}

func (tr *TableRenderer[T]) Top() {
	tr.SelectRow(0)
}

func (tr *TableRenderer[T]) Bottom() {
	tr.SelectRow(tr.Size() - 1)
}

func (tr *TableRenderer[T]) SelectRowByValue(f func(T) bool) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for idx, row := range tr.table.Rows {
		if f(row.Value) {
			tr.listRenderer.SetSelected(idx)
			return
		}
	}
}

func (tr *TableRenderer[T]) Up() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.listRenderer.Decrement()
}

func (tr *TableRenderer[T]) Down() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.listRenderer.Increment()
}

func (tr *TableRenderer[T]) Clear() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.table.clear()
	tr.listRenderer.ResetSize(0)
}

// Fills the table.
func (tr *TableRenderer[T]) Fill(rows [][]string, rowValues []T) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if (len(rows) > 0 && len(rows[0]) != len(tr.table.Title.Titles)) || len(rowValues) != len(rows) {
		log.Panicln("invalid row length")
	}

	tr.table.clear()
	for idx, row := range rows {
		tr.table.addTableRow(row, rowValues[idx])
	}

	tr.listRenderer.ResetSize(len(rows))
}

// Wrapper over RenderTable that passes default call backs.
func (tr *TableRenderer[T]) Render(w io.Writer) {
	tr.RenderTable(w, nil, nil)
}

// Renders the table to the Writer with the passed callbacks as optional.
// onSelected will be called with the selected row if passed.
// update will be called at every row if passed
func (tr *TableRenderer[T]) RenderTable(w io.Writer, onSelected func(int, *TableRow[T]) bool, update func(int, *TableRow[T])) {
	if update != nil {
		// if update function is not nil, we Lock for a write since the row might be updated during render
		tr.mu.Lock()
		defer tr.mu.Unlock()
	} else {
		// otherwise no updates are done so we RLock
		tr.mu.RLock()
		defer tr.mu.RUnlock()
	}

	// render col table
	tr.renderTitle(w)

	// for each element that should be shown in the list
	tr.listRenderer.forEach(func(idx int) {
		// get row
		currentRow := tr.table.Rows[idx]

		// set selected by checking if the selected is this row
		// but a call back can be passed to modify this
		currentRow.Selected = tr.listRenderer.selected == idx
		if currentRow.Selected && onSelected != nil {
			currentRow.Selected = onSelected(idx, currentRow)
		}

		if update != nil {
			update(idx, currentRow)
		}

		// render cols of this row
		var line string
		for i, col := range currentRow.Cols {
			proportion := tr.table.ColProportions[i]
			colSize := proportion * float64(tr.table.Width)
			colLine := Pad(col, int(math.Floor(colSize)))

			var style color.Style
			if currentRow.Selected {
				style = color.New(color.FgBlack, color.BgCyan)
			} else {
				style = tr.table.Title.Styles[i]
			}
			line += style.Sprint(colLine)
		}

		fmt.Fprintln(w, line)
	})
}

func (tr *TableRenderer[T]) renderTitle(w io.Writer) {
	var line string
	for i, title := range tr.table.Title.Titles {
		proportion := tr.table.ColProportions[i]
		colSize := proportion * float64(tr.table.Width)
		colLine := Pad(title, int(math.Floor(colSize)))

		line += colLine
	}

	s := color.Note.Style
	line = s.Sprint(line)
	fmt.Fprintln(w, line)
}

func (t *Table[T]) addTableRow(cols []string, value T) {
	tr := &TableRow[T]{
		Cols:  cols,
		Value: value,
	}

	t.Rows = append(t.Rows, tr)
}

func (t *Table[T]) clear() {
	t.Rows = make([]*TableRow[T], 0)
}
