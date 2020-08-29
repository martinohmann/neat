package table

import (
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/text"
)

// Option is a func for configuring a *Table.
type Option func(t *Table)

// WithPadding sets the horizontal padding between adjacent table cells.
// Defaults to 1.
func WithPadding(padding int) Option {
	return func(t *Table) {
		t.padding = padding
	}
}

// WithMargin sets the left and right margin of table rows. Defaults to 0.
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

// WithBorderMask controls the borders that should be displayed using a bit
// mask.
func WithBorderMask(mask BorderMask) Option {
	return func(t *Table) {
		t.borderMask = mask
	}
}

// WithBorderRunes sets the runes that should be rendered for the individual
// border elements, e.g. corners, vertical and horizontal lines and junctions.
// See the documentation of DefaultBorderRunes for an example. It is valid to
// pass a rune mapping for a subset of the required border runes. Missing runes
// will be filled with the corresponding runes from DefaultBorderRunes.
func WithBorderRunes(runes BorderRunes) Option {
	return func(t *Table) {
		t.borderRunes = runes
	}
}

// WithBorderStyle sets the style the should be applied to each border element.
// The default is to not apply any style.
func WithBorderStyle(style *style.Style) Option {
	return func(t *Table) {
		t.borderStyle = style
	}
}

// WithAlignment sets the default alignment for all table cells. Can be
// overridden per table column via WithColumnAlignment. The alignment is
// ignored for console.Renderable values passed to table.AddRow as these may
// have their own alignment rules. See the documentation of table.AddRow.
func WithAlignment(alignment text.Alignment) Option {
	return func(t *Table) {
		t.alignment = alignment
	}
}

// WithAlignment sets the default style for all table cells. Can be overridden
// per table column via WithColumnStyle. The style is ignored for
// console.Renderable values passed to table.AddRow as these may have their own
// style rules. See the documentation of table.AddRow.
func WithStyle(style *style.Style) Option {
	return func(t *Table) {
		t.style = style
	}
}

// WithWordWrap sets the default word wrapping behaviour for all table cells.
// If true, cells are word wrapped if width constraints make this necessary.
// Otherwise, cells are not word wrapped but may be truncated if they exceed
// the width constraints. Can be overridden per table column via
// WithColumnWordWrap. The word wrap config is ignored for console.Renderable
// values passed to table.AddRow as these may have their own word wrapping
// rules. See the documentation of table.AddRow.
func WithWordWrap(wrap bool) Option {
	return func(t *Table) {
		t.wordWrap = wrap
	}
}

// WithColumnAlignment configures the cell alignment per column. See
// documentation of WithAlignment.
func WithColumnAlignment(alignment ...text.Alignment) Option {
	return func(t *Table) {
		t.columnAlignment = alignment
	}
}

// WithColumnStyle configures the cell style per column. See
// documentation of WithStyle.
func WithColumnStyle(style ...*style.Style) Option {
	return func(t *Table) {
		t.columnStyle = style
	}
}

// WithWordWrap configures the cell word wrapping behaviour per column. See
// documentation of WithWordWrap.
func WithColumnWordWrap(wrap ...bool) Option {
	return func(t *Table) {
		t.columnWordWrap = wrap
	}
}
