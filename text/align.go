package text

import (
	"strings"
)

const (
	space    rune = ' '
	newline  rune = '\n'
	ellipsis rune = '…'
)

type Alignment int

const (
	AlignLeft Alignment = iota
	AlignRight
	AlignCenter
	AlignJustify
)

func Align(s string, width int, align Alignment) string {
	switch align {
	case AlignRight:
		return PadLeft(s, width)
	case AlignCenter:
		return PadCenter(s, width)
	case AlignJustify:
		return PadJustify(s, width)
	default:
		return PadRight(s, width)
	}
}

func PadLeft(s string, width int) string {
	return padLines(s, width, func(s string, padding int) string {
		return Spaces(padding) + s
	})
}

func PadRight(s string, width int) string {
	return padLines(s, width, func(s string, padding int) string {
		return s + Spaces(padding)
	})
}

func PadCenter(s string, width int) string {
	return padLines(s, width, func(s string, padding int) string {
		paddingLeft := int(float64(padding / 2))
		paddingRight := padding - paddingLeft

		return Spaces(paddingLeft) + s + Spaces(paddingRight)
	})
}

func PadJustify(s string, width int) string {
	return padLines(s, width, func(s string, padding int) string {
		if len(s) == 0 {
			return Spaces(padding)
		}

		words := strings.Split(s, string(space))

		if len(words) == 1 {
			return s + Spaces(padding)
		}

		// Calculate additional padding needed between adjacent words.
		wordPadding := int(float64(padding / (len(words) - 1)))

		var justified string

		for i, word := range words {
			// Last word receives remaining padding
			if i >= len(words)-1 && padding > wordPadding {
				wordPadding = padding
			}

			// First word is not padded
			if i > 0 {
				// +1 for existing space
				justified += Spaces(wordPadding + 1)
				padding -= wordPadding
			}

			justified += word
		}

		return justified
	})
}

func pad(s string, width int, fn func(string, int) string) string {
	padding := width - DisplayWidth(s)
	if padding > 0 {
		s = fn(s, padding)
	}

	return s
}

func padLines(s string, width int, fn func(string, int) string) string {
	if !IsMultiline(s) {
		return pad(s, width, fn)
	}

	lines := SplitLines(s)
	for i, line := range lines {
		lines[i] = pad(line, width, fn)
	}

	return JoinLines(lines)
}
