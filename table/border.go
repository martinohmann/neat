package table

// BorderMask controls which borders should be displayed.
type BorderMask int

// Has returns true if b has the bits of mask set.
func (b BorderMask) Has(mask BorderMask) bool {
	return b&mask == mask
}

// BorderMask options.
const (
	BorderNone BorderMask = 1 << iota
	BorderLeft
	BorderColumn
	BorderRight
	BorderTop
	BorderRow
	BorderBottom
	BorderSection

	BorderAllHorizontal = BorderTop | BorderRow | BorderBottom | BorderSection
	BorderAllVertical   = BorderLeft | BorderColumn | BorderRight
	BorderAll           = BorderAllHorizontal | BorderAllVertical
)

// BorderRune indicates the type of rune used for a border.
type BorderRune int

// BorderRune elements that are needed to draw a table with corners,
// intersections, row separators (horizontal) and column separators (vertical
// borders). Can be configured on a table via the WithBorderRunes table option.
// See DefaultBorderRunes for a mapping of these constants to actual runes.
const (
	BorderRuneHorizontal                BorderRune = iota // ─
	BorderRuneVertical                                    // │
	BorderRuneCornerTopLeft                               // ┌
	BorderRuneCornerTopRight                              // ┐
	BorderRuneCornerBottomLeft                            // └
	BorderRuneCornerBottomRight                           // └
	BorderRuneIntersectionTop                             // ┬
	BorderRuneIntersectionBottom                          // ┴
	BorderRuneIntersectionLeft                            // ├
	BorderRuneIntersectionRight                           // ┤
	BorderRuneIntersectionCenter                          // ┼
	BorderRuneSectionHorizontal                           // ═
	BorderRuneSectionIntersectionLeft                     // ╞
	BorderRuneSectionIntersectionRight                    // ╡
	BorderRuneSectionIntersectionCenter                   // ╪
)

// BorderRunes is a map of the BorderRune type to the actual rune that should
// be displayed.
type BorderRunes map[BorderRune]rune

// DefaultBorderRunes are the runes that will be used to draw table borders if
// not explicitly overridden via table options. This is an exported variable to
// allow overriding table borders globally.
var DefaultBorderRunes = BorderRunes{
	BorderRuneHorizontal:                '─',
	BorderRuneVertical:                  '│',
	BorderRuneCornerTopLeft:             '┌',
	BorderRuneCornerTopRight:            '┐',
	BorderRuneCornerBottomLeft:          '└',
	BorderRuneCornerBottomRight:         '┘',
	BorderRuneIntersectionTop:           '┬',
	BorderRuneIntersectionBottom:        '┴',
	BorderRuneIntersectionLeft:          '├',
	BorderRuneIntersectionRight:         '┤',
	BorderRuneIntersectionCenter:        '┼',
	BorderRuneSectionHorizontal:         '═',
	BorderRuneSectionIntersectionLeft:   '╞',
	BorderRuneSectionIntersectionRight:  '╡',
	BorderRuneSectionIntersectionCenter: '╪',
}
