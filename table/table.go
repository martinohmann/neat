package table

import (
	"fmt"
	"io"

	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/internal/util"
	"github.com/martinohmann/neat/measure"
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/text"
	runewidth "github.com/mattn/go-runewidth"
)

const defaultMaxWidth = 80

// Table can render properly aligned columns and rows of information.
type Table struct {
	out      io.Writer
	padding  int
	margin   int
	maxWidth int

	borderMask  BorderMask
	borderRunes BorderRunes
	borderStyle *style.Style

	// global cell attributes
	alignment text.Alignment
	style     *style.Style
	wordWrap  bool

	// per-column cell attributes
	columnAlignment []text.Alignment
	columnStyle     []*style.Style
	columnWordWrap  []bool

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

	t.applyOptions(opts)

	return t
}

func (t *Table) applyOptions(opts []Option) {
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

	if t.borderRunes == nil {
		t.borderRunes = DefaultBorderRunes
	} else {
		for k := range DefaultBorderRunes {
			if _, ok := t.borderRunes[k]; !ok {
				t.borderRunes[k] = DefaultBorderRunes[k]
			}
		}
	}
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
//
// Columns implementing console.Renderable are NOT formatted using the cell and
// column specific options (e.g. style, alignment, word wrap) configured via
// the table.With* and table.WithColumn* option funcs. This allows users to add
// custom cell behaviour if needed.
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

// calculateSpacing calculates the width and height occupied by spacing like
// padding, margin and borders.
func (t *Table) calculateSpacing() (width, height int) {
	borderWidth := 0
	paddingWidth := (len(t.cols) - 1) * t.padding
	marginWidth := 2 * t.margin

	width = marginWidth + paddingWidth

	for _, r := range t.borderRunes {
		borderWidth = util.MaxInt(borderWidth, runewidth.RuneWidth(r))
	}

	if t.borderMask.Has(BorderColumn) {
		// If we have vertical borders we need to have double the padding: left
		// and right from the border.
		width += paddingWidth + ((len(t.cols) - 1) * borderWidth)
	}

	if t.borderMask.Has(BorderLeft) {
		width += t.padding + borderWidth
	}

	if t.borderMask.Has(BorderRight) {
		width += t.padding + borderWidth
	}

	if t.borderMask.Has(BorderTop) {
		height++
	}

	if t.borderMask.Has(BorderRow) {
		height += len(t.rows) - 1
	}

	if t.borderMask.Has(BorderBottom) {
		height++
	}

	return width, height
}

// Render renders the table to the underlying io.Writer. Returns the number of
// lines rendered and an error if rendering failed.
func (t *Table) Render() (int, error) {
	if len(t.cols) == 0 {
		return 0, nil
	}

	spacingWidth, spacingHeight := t.calculateSpacing()
	availWidth := t.maxWidth - spacingWidth

	measures := t.measureColumns(availWidth)
	columnWidths := measure.Sum(measures...)

	totalWidth := spacingWidth + columnWidths.Maximum

	tb := newTableBuilder(t, measures)

	tb.Grow((len(t.rows) + spacingHeight) * (totalWidth + 1))

	if t.borderMask.Has(BorderTop) {
		tb.writeBorderLine(BorderRuneCornerTopLeft, BorderRuneIntersectionTop, BorderRuneCornerTopRight)
	}

	for i, row := range t.rows {
		cellLines, maxCellHeight := renderRowCells(row, measures)

		if maxCellHeight > 1 {
			// Grow buffer to have enough space to hold all lines of the table
			// row without the need to reallocate between writing the lines of
			// the row.
			tb.Grow((maxCellHeight - 1) * (totalWidth + 1))
		}

		tb.writeRowCells(cellLines, maxCellHeight)

		if t.borderMask.Has(BorderRow) && i < len(t.rows)-1 {
			tb.writeBorderLine(BorderRuneIntersectionLeft, BorderRuneIntersectionCenter, BorderRuneIntersectionRight)
		}
	}

	if t.borderMask.Has(BorderBottom) {
		tb.writeBorderLine(BorderRuneCornerBottomLeft, BorderRuneIntersectionBottom, BorderRuneCornerBottomRight)
	}

	return tb.render()
}

// renderRowCells renders all cells of the row to produce the lines that we
// need to calculate the row height. Returns a slice of slices of lines for
// each column in the row and the maximum column height.
func renderRowCells(row *tableRow, measures []measure.Measurement) (cellLines [][]string, maxHeight int) {
	cellLines = make([][]string, len(row.cells))

	for colIdx, cell := range row.cells {
		rendered := cell.Render(measures[colIdx].Maximum)
		lines := text.SplitLines(rendered)

		maxHeight = util.MaxInt(maxHeight, len(lines))

		cellLines[colIdx] = lines
	}

	return cellLines, maxHeight
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

func (t *Table) makeRenderable(v interface{}, colIdx int) console.Renderable {
	if r, ok := v.(console.Renderable); ok {
		return r
	}

	return t.makeTextRenderable(v, colIdx)
}

func (t *Table) makeTextRenderable(v interface{}, colIdx int) console.Renderable {
	r := text.Text{
		Alignment: t.alignment,
		Style:     t.style,
		Text:      fmt.Sprint(v),
		WordWrap:  t.wordWrap,
	}

	if colIdx < len(t.columnAlignment) {
		r.Alignment = t.columnAlignment[colIdx]
	}

	if colIdx < len(t.columnStyle) {
		r.Style = t.columnStyle[colIdx]
	}

	if colIdx < len(t.columnWordWrap) {
		r.WordWrap = t.columnWordWrap[colIdx]
	}

	return r
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
