package ui

import (
	"fmt"
	"io"
	"log"
	"math"

	"github.com/gookit/color"
)

type (
	TableRenderer struct {
		Table        *Table
		ListRenderer *ListRenderer
	}

	Table struct {
		Title          *TableTitle
		Rows           []*TableRow
		ColProportions []float64

		Width  int
		Height int
	}

	TableTitle struct {
		Titles []string
	}

	TableRow struct {
		Cols     []string
		Selected bool
	}
)

func NewTableRenderer() *TableRenderer {
	return &TableRenderer{}
}

func (tr *TableRenderer) InitTable(width int, height int, titles []string, colProportions []float64) {
	if len(titles) != len(colProportions) {
		log.Panicln("the number of titles and col proportions should be the same")
	}
	tr.ListRenderer = newListRenderer(0, height-1, 0)
	tr.Table = &Table{
		Title: &TableTitle{
			Titles: titles,
		},
		ColProportions: colProportions,
		Rows:           make([]*TableRow, 0),

		Width:  width,
		Height: height,
	}
}

func (t *Table) AddTableRow(cols []string) {
	tr := &TableRow{
		Cols: cols,
	}

	t.Rows = append(t.Rows, tr)
}

func (t *Table) ClearTable() {
	t.Rows = make([]*TableRow, 0)
}

func (tr *TableRenderer) GetSelectedRowIndex() int {
	return tr.ListRenderer.selected
}

func (tr *TableRenderer) SetSelectedRow(idx int) {
	tr.ListRenderer.setSelected(idx)
}

func (tr *TableRenderer) Up() {
	tr.ListRenderer.decrement()
}

func (tr *TableRenderer) Down() {
	tr.ListRenderer.increment()
}

func (tr *TableRenderer) RenderWithSelectCallBack(w io.Writer, onSelected func(int, *TableRow) bool) {
	tr.render(w, onSelected)
}

func (tr *TableRenderer) Render(w io.Writer) {
	tr.render(w, func(i int, tr *TableRow) bool { return true })
}

func (tr *TableRenderer) render(w io.Writer, onSelected func(int, *TableRow) bool) {
	tr.RenderTitle(w)

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

func (tr *TableRenderer) RenderTitle(w io.Writer) {
	var line string
	for i, title := range tr.Table.Title.Titles {
		proportion := tr.Table.ColProportions[i]
		colSize := proportion * float64(tr.Table.Width)
		colLine := withSpacePadding(title, int(math.Ceil(colSize)))

		line += colLine
	}

	fmt.Fprintln(w, line)
}

func (tr *TableRenderer) FillTable(rows [][]string) {
	if len(rows) > 0 && len(rows[0]) != len(tr.Table.Title.Titles) {
		log.Panicln("invalid row length")
	}

	tr.Table.ClearTable()
	for _, row := range rows {
		tr.Table.AddTableRow(row)
	}
	tr.ListRenderer.resetSize(len(rows))
}
