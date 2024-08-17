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

func (tr *TableRenderer[T]) InitTable(width int, height int, titles []string, colProportions []float64) {
	if len(titles) != len(colProportions) {
		log.Panicln("the number of titles and col proportions should be the same")
	}
	tr.listRenderer = NewListRenderer(0, height-1, 0)
	tr.table = &Table[T]{
		Title: &TableTitle{
			Titles: titles,
		},
		ColProportions: colProportions,
		Rows:           make([]*TableRow[T], 0),

		Width:  width,
		Height: height,
	}
}

func (tr *TableRenderer[T]) GetTableSize() int {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return len(tr.table.Rows)
}

func (tr *TableRenderer[T]) GetSelectedRow() (idx int, value *T) {
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

func (tr *TableRenderer[T]) ClearTable() {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.table.clear()
	tr.listRenderer.ResetSize(0)
}

func (tr *TableRenderer[T]) FillTable(rows [][]string, rowValues []T) {
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

	tr.listRenderer.forEach(func(idx int) {
		currentRow := tr.table.Rows[idx]
		currentRow.Selected = tr.listRenderer.selected == idx
		if currentRow.Selected {
			currentRow.Selected = onSelected(idx, currentRow)
		}

		var line string
		for i, col := range currentRow.Cols {
			proportion := tr.table.ColProportions[i]
			colSize := proportion * float64(tr.table.Width)
			colLine := WithSpaces(col, int(math.Ceil(colSize)))
			line += colLine
		}

		var c color.Style
		if currentRow.Selected {
			c = color.New(color.Black, color.BgCyan)
		} else {
			c = color.New(color.White)
		}

		line = c.Sprint(line)
		fmt.Fprintln(w, line)
	})
}

func (tr *TableRenderer[T]) renderTitle(w io.Writer) {
	var line string
	for i, title := range tr.table.Title.Titles {
		proportion := tr.table.ColProportions[i]
		colSize := proportion * float64(tr.table.Width)
		colLine := WithSpaces(title, int(math.Ceil(colSize)))

		line += colLine
	}

	s := color.New(color.Blue)
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
