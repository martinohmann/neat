package table

import (
	"fmt"
	"strings"

	"github.com/martinohmann/neat/measure"
	"github.com/martinohmann/neat/text"
)

// tableBuilder builds the string representation of a table and writes it to
// the underlying io.Writer of the wrapped *Table.
type tableBuilder struct {
	strings.Builder
	*Table

	measures      []measure.Measurement
	marginSpaces  string
	paddingSpaces string
	lines         int
}

func newTableBuilder(t *Table, measures []measure.Measurement) *tableBuilder {
	return &tableBuilder{
		Table:         t,
		measures:      measures,
		marginSpaces:  text.Spaces(t.margin),
		paddingSpaces: text.Spaces(t.padding),
	}
}

func (tb *tableBuilder) writeBorderString(s string) {
	if tb.borderStyle != nil {
		s = tb.borderStyle.Sprint(s)
	}

	tb.WriteString(s)
}

func (tb *tableBuilder) writeBorderRuneN(r BorderRune, n int) {
	tb.writeBorderString(strings.Repeat(string(tb.borderRunes[r]), n))
}

func (tb *tableBuilder) writeBorderRune(r BorderRune) {
	tb.writeBorderString(string(tb.borderRunes[r]))
}

func (tb *tableBuilder) writePaddingSpaces() { tb.WriteString(tb.paddingSpaces) }

func (tb *tableBuilder) writeMarginSpaces() { tb.WriteString(tb.marginSpaces) }

func (tb *tableBuilder) writeNewline() {
	tb.WriteRune('\n')
	tb.lines++
}

func (tb *tableBuilder) writeRowCells(cellLines [][]string, maxCellHeight int) {
	// Write all cells of the current row to the buffer and handle multiple
	// lines.
	for lineNum := 0; lineNum < maxCellHeight; lineNum++ {
		tb.writeMarginSpaces()

		if tb.borderMask.Has(BorderLeft) {
			tb.writeBorderRune(BorderRuneVertical)
			tb.writePaddingSpaces()
		}

		for colIdx, lines := range cellLines {
			if lineNum < len(lines) {
				tb.WriteString(lines[lineNum])
			} else {
				tb.WriteString(text.Spaces(tb.measures[colIdx].Maximum))
			}

			// Insert padding after each column except the last one.
			if colIdx < len(cellLines)-1 {
				tb.writePaddingSpaces()

				if tb.borderMask.Has(BorderColumn) {
					tb.writeBorderRune(BorderRuneVertical)
					tb.writePaddingSpaces()
				}
			}
		}

		if tb.borderMask.Has(BorderRight) {
			tb.writePaddingSpaces()
			tb.writeBorderRune(BorderRuneVertical)
		}

		tb.writeMarginSpaces()
		tb.writeNewline()
	}
}

func (tb *tableBuilder) writeBorderLine(left, junction, right BorderRune) {
	tb.writeMarginSpaces()

	if tb.borderMask.Has(BorderLeft) {
		tb.writeBorderRune(left)
		tb.writeBorderRuneN(BorderRuneHorizontal, tb.padding)
	}

	for i, measure := range tb.measures {
		line := strings.Repeat(string(tb.borderRunes[BorderRuneHorizontal]), measure.Maximum)

		tb.writeBorderString(line)

		if i < len(tb.measures)-1 {
			tb.writeBorderRuneN(BorderRuneHorizontal, tb.padding)

			if tb.borderMask.Has(BorderColumn) {
				tb.writeBorderRune(junction)
				tb.writeBorderRuneN(BorderRuneHorizontal, tb.padding)
			}
		}
	}

	if tb.borderMask.Has(BorderRight) {
		tb.writeBorderRuneN(BorderRuneHorizontal, tb.padding)
		tb.writeBorderRune(right)
	}

	tb.writeMarginSpaces()
	tb.writeNewline()
}

func (tb *tableBuilder) render() (int, error) {
	_, err := fmt.Fprint(tb.out, tb.String())

	return tb.lines, err
}
