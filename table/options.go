package table

import "github.com/martinohmann/neat/text"

// Option is a func for configuring a *Table.
type Option func(t *Table)

// WithPadding sets the horizontal padding between adjacent table cells.
// Defaults to 1.
func WithPadding(padding int) Option {
	return func(t *Table) {
		t.padding = padding
	}
}

// WidthMargin sets the left and right margin of table rows. Defaults to 0.
func WithMargin(margin int) Option {
	return func(t *Table) {
		t.margin = margin
	}
}

// WithMaxWidth sets the maximum table width. If maxWidth is <= 0, maxWidth is
// inferred from the table's underlying io.Writer if it is a
// console.FileWriter, otherwise a default of 80 is used.
func WithMaxWidth(maxWidth int) Option {
	return func(t *Table) {
		t.maxWidth = maxWidth
	}
}

func WithAlignment(alignments ...text.Alignment) Option {
	return func(t *Table) {
		t.alignments = alignments
	}
}
