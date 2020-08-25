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

const defaultMaxWidth = 80

// Table can render properly aligned columns and rows of information.
type Table struct {
	out      io.Writer
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

// New creates a new *Table which will be rendered to the provided io.Writer
// using opts.
func New(out io.Writer, opts ...Option) *Table {
	t := &Table{
		out:     out,
		padding: 1,
		margin:  0,
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

	if t.maxWidth <= 0 {
		if fw, ok := t.out.(console.FileWriter); ok {
			t.maxWidth = console.TerminalWidth(fw)
		} else {
			t.maxWidth = defaultMaxWidth
		}
	}

	return t
}

func (t *Table) mustFitWidth(cols []interface{}) {
	if len(t.rows) == 0 || len(cols) == len(t.rows[0].cells) {
		return
	}

	panic(fmt.Sprintf("expected %d columns, got %d", len(t.rows[0].cells), len(cols)))
}

// Reset resets the table by clearing all rows. This is useful for creating
// multiple tables with the same options.
func (t *Table) Reset() *Table {
	t.rows = nil
	t.cols = nil
	return t
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

// Render renders the table to the underlying io.Writer. Returns the number of
// lines rendered and an error if rendering failed.
func (t *Table) Render() (nlines int, err error) {
	if len(t.cols) == 0 {
		return 0, nil
	}

	paddingWidth := (len(t.cols) - 1) * t.padding
	marginWidth := 2 * t.margin
	availWidth := t.maxWidth - marginWidth - paddingWidth

	measures := t.measureColumns(availWidth)

	columnWidths := measure.Sum(measures...)

	totalWidth := marginWidth + paddingWidth + columnWidths.Maximum

	var sb strings.Builder

	sb.Grow(len(t.rows) * (totalWidth + 1))

	marginSpaces := text.Spaces(t.margin)
	paddingSpaces := text.Spaces(t.padding)

	for _, row := range t.rows {
		// Render all cells of the row to produce the lines that we need to
		// calculate the row height.
		rowHeight := 0
		cellLines := make([][]string, len(row.cells))

		for colIdx, cell := range row.cells {
			rendered := cell.Render(measures[colIdx].Maximum)
			lines := text.SplitLines(rendered)

			rowHeight = util.MaxInt(rowHeight, len(lines))

			cellLines[colIdx] = lines
		}

		if rowHeight > 1 {
			// Grow string buffer to have enough space to hold all lines of the
			// table row without the need to reallocate between writing rows.
			sb.Grow((rowHeight - 1) * (totalWidth + 1))
		}

		// Write all cells of the current row to the buffer and handle multiple
		// cells.
		for lineNum := 0; lineNum < rowHeight; lineNum++ {
			sb.WriteString(marginSpaces)

			for colIdx := range row.cells {
				lines := cellLines[colIdx]
				if lineNum < len(lines) {
					sb.WriteString(lines[lineNum])
				} else {
					sb.WriteString(text.Spaces(measures[colIdx].Maximum))
				}

				// Insert padding after each column except the last one.
				if colIdx < len(row.cells)-1 {
					sb.WriteString(paddingSpaces)
				}
			}

			sb.WriteString(marginSpaces)
			sb.WriteRune('\n')

			nlines++
		}
	}

	_, err = fmt.Fprint(t.out, sb.String())

	return nlines, err
}

func (t *Table) measureColumns(availWidth int) []measure.Measurement {
	measures := make([]measure.Measurement, len(t.cols))

	for i, col := range t.cols {
		measures[i] = col.measure(availWidth)
	}

	requested := measure.Sum(measures...)

	// Best case: columns fit nicely into the available space.
	if requested.Maximum <= availWidth {
		return measures
	}

	// Second best case: the optimal column widths overflow, but the minimum
	// requested space fits into the available space.
	if requested.Minimum <= availWidth {
		return t.overflowColumns(measures, availWidth)
	}

	// Worst case: we need to truncate to fit columns.
	return t.truncateColumns(measures, availWidth)
}

// overflowColumns tries to allocate space for all columns first that request
// less than the maximum width for each column if the available space is
// distributed evenly. Columns that are still unallocated after that will receive
// their requested minimum in the worst case, even if this exceeds the fair
// column share.
func (t *Table) overflowColumns(measures []measure.Measurement, availWidth int) []measure.Measurement {
	fairColWidth := int(float64(availWidth) / float64(len(t.cols)))

	remainingCols := len(t.cols)
	unallocated := make(map[int]struct{})

	for i, m := range measures {
		if m.Maximum > fairColWidth {
			unallocated[i] = struct{}{}
			continue
		}

		availWidth -= m.Maximum
		remainingCols--
	}

	for i, m := range measures {
		if _, ok := unallocated[i]; !ok {
			continue
		}

		fairColWidth = int(float64(availWidth) / float64(remainingCols))

		width := util.MinInt(availWidth, util.MaxInt(m.Minimum, fairColWidth))

		measures[i] = measure.NewMeasurement(util.MinInt(m.Minimum, width), width)
		availWidth -= width
		remainingCols--
	}

	return measures
}

// truncateColumns works similar to overflowColumns but aggressively truncates
// columns if the available with is not enough to display them all. This will
// truncate columns to less than their requested minimum if there is no other
// option.
func (t *Table) truncateColumns(measures []measure.Measurement, availWidth int) []measure.Measurement {
	fairColWidth := util.MaxInt(0, int(float64(availWidth)/float64(len(t.cols))))

	remainingCols := len(t.cols)
	unallocated := make(map[int]struct{})

	for i, m := range measures {
		if m.Minimum > fairColWidth {
			unallocated[i] = struct{}{}
			continue
		}

		measures[i].Maximum = m.Minimum
		availWidth -= m.Minimum
		remainingCols--
	}

	for i := range measures {
		if _, ok := unallocated[i]; !ok {
			continue
		}

		width := util.MaxInt(0, int(float64(availWidth)/float64(remainingCols)))
		measures[i] = measure.NewMeasurement(width, width)
		availWidth -= width
		remainingCols--
	}

	return measures
}

func (t *Table) makeCells(cols []interface{}) []console.Renderable {
	cells := make([]console.Renderable, len(cols))

	for i, col := range cols {
		cells[i] = t.makeRenderable(col, i)
	}

	return cells
}

func (t *Table) columnAlignment(colIdx int) *text.Alignment {
	if len(t.alignments) > colIdx {
		return &t.alignments[colIdx]
	}

	return nil
}

func (t *Table) makeRenderable(v interface{}, colIdx int) console.Renderable {
	alignment := t.columnAlignment(colIdx)
	// FIXME: the way we handle column alignment here needs to be improved. For
	// now we leave it like this.
	if alignment == nil {
		switch t := v.(type) {
		case console.Renderable:
			return t
		default:
			return text.Text{Text: fmt.Sprint(t)}
		}
	} else {
		switch t := v.(type) {
		case *text.Text:
			tv := *t
			tv.Alignment = *alignment
			return tv
		case text.Text:
			t.Alignment = *alignment
			return t
		case console.Renderable:
			return t
		default:
			return text.Text{Text: fmt.Sprint(t), Alignment: *alignment}
		}
	}
}

type tableRow struct {
	cells []console.Renderable
}

type tableCol struct {
	cells []console.Renderable
}

func (c *tableCol) measure(maxWidth int) measure.Measurement {
	var m measure.Measurement

	for _, cell := range c.cells {
		cm := cell.Measure(maxWidth)

		m = measure.NewMeasurement(
			util.MaxInt(m.Minimum, cm.Minimum),
			util.MaxInt(m.Maximum, cm.Maximum),
		)
	}

	return m
}
