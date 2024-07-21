package ui

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
		Table        *Table[T]
		ListRenderer *ListRenderer
		mu           *sync.RWMutex
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
	}

	TableRow[T any] struct {
		Value    T
		Cols     []string
		Selected bool
	}
)

func NewTableRenderer[T any]() *TableRenderer[T] {
	return &TableRenderer[T]{
		mu: &sync.RWMutex{},
	}
}

func (tr *TableRenderer[T]) InitTable(width int, height int, titles []string, colProportions []float64) {
	if len(titles) != len(colProportions) {
		log.Panicln("the number of titles and col proportions should be the same")
	}
	tr.ListRenderer = newListRenderer(0, height-1, 0)
	tr.Table = &Table[T]{
		Title: &TableTitle{
			Titles: titles,
		},
		ColProportions: colProportions,
		Rows:           make([]*TableRow[T], 0),

		Width:  width,
		Height: height,
	}
}

func (tr TableRenderer[T]) GetTableSize() int {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return len(tr.Table.Rows)
}

func (tr *TableRenderer[T]) GetSelectedRow() (idx int, value *T) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	if len(tr.Table.Rows) == 0 {
		return 0, nil
	}
	return tr.ListRenderer.selected, &tr.Table.Rows[tr.ListRenderer.selected].Value
}

func (tr *TableRenderer[T]) SetSelectedRow(idx int) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.ListRenderer.setSelected(idx)
}

// TODO: RENAME
func (tr *TableRenderer[T]) SetSelectedRowByValue(f func(T) bool) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	for idx, row := range tr.Table.Rows {
		if f(row.Value) {
			tr.ListRenderer.setSelected(idx)
			return
		}
	}
}

func (tr *TableRenderer[T]) Up() {
	tr.ListRenderer.decrement()
}

func (tr *TableRenderer[T]) Down() {
	tr.ListRenderer.increment()
}

func (tr *TableRenderer[T]) FillTable(rows [][]string, rowValues []T) {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if (len(rows) > 0 && len(rows[0]) != len(tr.Table.Title.Titles)) || len(rowValues) != len(rows) {
		log.Panicln("invalid row length")
	}

	tr.Table.Clear()
	for idx, row := range rows {
		tr.Table.addTableRow(row, rowValues[idx])
	}

	tr.ListRenderer.resetSize(len(rows))
}

func (tr *TableRenderer[T]) RenderWithSelectCallBack(w io.Writer, onSelected func(int, *TableRow[T]) bool) {
	tr.render(w, onSelected)
}

func (tr *TableRenderer[T]) Render(w io.Writer) {
	tr.render(w, func(i int, tr *TableRow[T]) bool { return true })
}

func (tr *TableRenderer[T]) render(w io.Writer, onSelected func(int, *TableRow[T]) bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	tr.renderTitle(w)

	tr.ListRenderer.forEach(func(idx int) {
		currentRow := tr.Table.Rows[idx]
		currentRow.Selected = tr.ListRenderer.selected == idx
		if currentRow.Selected {
			currentRow.Selected = onSelected(idx, currentRow)
		}

		var line string
		for i, col := range currentRow.Cols {
			proportion := tr.Table.ColProportions[i]
			colSize := proportion * float64(tr.Table.Width)
			colLine := withSpacePadding(col, int(math.Ceil(colSize)))
			line += colLine
		}

		var c color.Style
		if currentRow.Selected {
			c = color.New(color.Black, color.BgCyan)
		} else {
			c = color.New(color.Blue)
		}
		line = c.Sprint(line)

		fmt.Fprintln(w, line)
	})
}

func (tr *TableRenderer[T]) renderTitle(w io.Writer) {
	var line string
	for i, title := range tr.Table.Title.Titles {
		proportion := tr.Table.ColProportions[i]
		colSize := proportion * float64(tr.Table.Width)
		colLine := withSpacePadding(title, int(math.Ceil(colSize)))

		line += colLine
	}

	fmt.Fprintln(w, line)
}

func (t *Table[T]) addTableRow(cols []string, value T) {
	tr := &TableRow[T]{
		Cols:  cols,
		Value: value,
	}

	t.Rows = append(t.Rows, tr)
}

func (t *Table[T]) Clear() {
	t.Rows = make([]*TableRow[T], 0)
}
