package table

import (
	"fmt"
	"io"
	"strings"

	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/internal/util"
	"github.com/martinohmann/neat/measure"
	"github.com/martinohmann/neat/text"
)

// Table can render properly aligned columns and rows of information.
type Table struct {
	padding  int
	margin   int
	maxWidth int

	alignments []text.Alignment

	// rows and cols contain pointers to exactly the same cells, only in
	// different orientation. This helps to simplify calculation of column
	// widths and row heights.
	rows []*tableRow
	cols []*tableCol
}

// New creates a new *Table of maxWidth using opts.
func New(maxWidth int, opts ...Option) *Table {
	t := &Table{
		maxWidth: maxWidth,
		padding:  1,
		margin:   0,
	}

	for _, option := range opts {
		option(t)
	}

	if t.padding < 0 {
		t.padding = 0
	}

	if t.margin < 0 {
		t.margin = 0
	}

	return t
}

func (t *Table) mustFitWidth(cols []interface{}) {
	if len(t.rows) == 0 || len(cols) == len(t.rows[0].cells) {
		return
	}

	panic(fmt.Sprintf("expected %d columns, got %d", len(t.rows[0].cells), len(cols)))
}

// AddRow adds a row to the table. Panics if the number of columns does not
// align with the number of columns of already existing table rows.
func (t *Table) AddRow(columns ...interface{}) *Table {
	t.mustFitWidth(columns)

	cells := t.makeCells(columns)

	t.rows = append(t.rows, &tableRow{cells: cells})

	if t.cols == nil {
		t.cols = make([]*tableCol, len(cells))
	}

	for i, cell := range cells {
		if t.cols[i] == nil {
			t.cols[i] = &tableCol{}
		}

		t.cols[i].cells = append(t.cols[i].cells, cell)
	}

	return t
}

// Render renders the table and writes the result to w. Returns the number of
// lines rendered and an error if rendering failed.
func (t *Table) Render(w io.Writer) (nlines int, err error) {
	if len(t.cols) == 0 {
		return 0, nil
	}

	maxWidth := t.maxWidth

	dynamicCols := make([]int, 0)
	dynamicMinWidth := 0
	paddingWidth := (len(t.cols) - 1) * t.padding
	marginWidth := 2 * t.margin
	staticWidth := paddingWidth + marginWidth

	// Calculate the column widths and mark dynamic columns.
	for i, col := range t.cols {
		col.recalculateWidth(maxWidth)

		if col.width.Maximum >= maxWidth {
			dynamicCols = append(dynamicCols, i)
			dynamicMinWidth += col.width.Minimum
		} else {
			staticWidth += col.width.Maximum
		}
	}

	remainingWidth := maxWidth - staticWidth

	if remainingWidth < dynamicMinWidth {
		staticWidth = paddingWidth + marginWidth
		availWidth := maxWidth - staticWidth
		maxColWidth := int(float64(availWidth) / float64(len(t.cols)))

		for i, col := range t.cols {
			if containsInt(dynamicCols, i) {
				continue
			}

			col.recalculateWidth(maxColWidth)
			staticWidth += col.width.Maximum
		}

		remainingWidth = maxWidth - staticWidth
	}

	remainingWidth = util.MaxInt(0, remainingWidth)

	if len(dynamicCols) > 0 {
		dynamicColWidth := int(float64(remainingWidth) / float64(len(dynamicCols)))

		for _, idx := range dynamicCols {
			t.cols[idx].recalculateWidth(dynamicColWidth)
		}
	}

	var sb strings.Builder
	sb.Grow((maxWidth + 1) * len(t.rows))

	lineCount := 0

	marginSpaces := text.Spaces(t.margin)
	paddingSpaces := text.Spaces(t.padding)

	for i, row := range t.rows {
		// Render all cells of the row to produce the lines that we need to
		// calculate the row height.
		for j := 0; j < len(row.cells); j++ {
			col := t.cols[j]
			cell := col.cells[i]

			if cell.lines == nil {
				rendered := cell.Render(col.width.Maximum)

				cell.lines = text.SplitLines(rendered)
			}
		}

		rowHeight := row.height()

		if rowHeight > 1 {
			sb.Grow(maxWidth + 1)
		}

		lineCount += rowHeight

		for lineNum := 0; lineNum < rowHeight; lineNum++ {
			sb.WriteString(marginSpaces)

			for j := 0; j < len(row.cells); j++ {
				col := t.cols[j]
				cell := col.cells[i]

				if lineNum < len(cell.lines) {
					sb.WriteString(cell.lines[lineNum])
				} else {
					sb.WriteString(text.Spaces(col.width.Maximum))
				}

				// Insert padding after each column except the last one.
				if j < len(row.cells)-1 {
					sb.WriteString(paddingSpaces)
				}
			}

			sb.WriteString(marginSpaces)
			sb.WriteRune('\n')
		}
	}

	_, err = fmt.Fprint(w, sb.String())

	return lineCount, err
}

func containsInt(haystack []int, needle int) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func (t *Table) makeCells(cols []interface{}) []*tableCell {
	cells := make([]*tableCell, len(cols))

	for i, col := range cols {
		cells[i] = &tableCell{Renderable: t.makeRenderable(col, i)}
	}

	return cells
}

func (t *Table) columnAlignment(colIdx int) text.Alignment {
	if len(t.alignments) > colIdx {
		return t.alignments[colIdx]
	}

	return text.AlignLeft
}

func (t *Table) makeRenderable(v interface{}, colIdx int) console.Renderable {
	alignment := t.columnAlignment(colIdx)

	switch t := v.(type) {
	case *text.Text:
		tv := *t
		tv.Alignment = alignment
		return tv
	case text.Text:
		t.Alignment = alignment
		return t
	case console.Renderable:
		return t
	case string:
		return text.Text{Text: t, Alignment: alignment}
	case fmt.Stringer:
		return text.Text{Text: t.String(), Alignment: alignment}
	default:
		return text.Text{Text: fmt.Sprint(t), Alignment: alignment}
	}
}

type tableRow struct {
	cells []*tableCell
}

// height returns the height of the tallest cell in the row. If called before
// populating the lines slice of all row's cells this will return an inaccurate
// height.
func (r *tableRow) height() (h int) {
	for _, cell := range r.cells {
		h = util.MaxInt(h, len(cell.lines))
	}

	return h
}

type tableCell struct {
	console.Renderable
	// Populated by table.Render.
	lines []string
}

type tableCol struct {
	cells []*tableCell
	// Populated by recalculateWidth.
	width measure.Measurement
}

func (c *tableCol) recalculateWidth(maxWidth int) {
	var m measure.Measurement

	for _, cell := range c.cells {
		cm := cell.Measure(maxWidth)

		m = measure.NewMeasurement(
			util.MaxInt(m.Minimum, cm.Minimum),
			util.MaxInt(m.Maximum, cm.Maximum),
		)
	}

	c.width = m
}
