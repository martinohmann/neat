package text

import (
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/martinohmann/neat/internal/util"
	runewidth "github.com/mattn/go-runewidth"
)

// Truncate truncates s to a maximum width. The truncated string includes an
// ellipsis as the last rune. If the display width of s is shorter than width,
// it is returned unaltered. For truncate to work as expected s must not
// contain newlines.
func Truncate(s string, width int) string {
	if width == 0 {
		return ""
	}

	if width > 0 && displayWidth(s) > width {
		s = s[:width-1] + string(ellipsis)
	}

	return s
}

// DisplayWidth returns the display width of s. If s is a multiline string this
// returns the display width of the longest line.
func DisplayWidth(s string) int {
	if IsMultiline(s) {
		return MaxDisplayWidth(SplitLines(s))
	}

	return displayWidth(s)
}

func displayWidth(s string) int {
	s = stripansi.Strip(s)

	return runewidth.StringWidth(s)
}

// MaxDisplayWidth returns the display width of the longest line in the lines
// slice.
func MaxDisplayWidth(lines []string) (max int) {
	for _, line := range lines {
		max = util.MaxInt(max, displayWidth(line))
	}

	return max
}

// IsMultiline returns true if s contains newlines.
func IsMultiline(s string) bool {
	return strings.Contains(s, string(newline))
}

// CountLines returns the number of lines in s.
func CountLines(s string) int {
	return strings.Count(s, string(newline)) + 1
}

// SplitLines splits s on newline characters and returns a slice of lines.
func SplitLines(s string) []string {
	return strings.Split(s, string(newline))
}

// JoinLines joins lines together with newline characters.
func JoinLines(lines []string) string {
	return strings.Join(lines, string(newline))
}

// Spaces returns a string which consist only of the specified number of
// spaces. Will panic is num is negative.
func Spaces(num int) string {
	return strings.Repeat(string(space), num)
}

// This is based on the naive implementation found here:
// https://www.rosettacode.org/wiki/Word_wrap#Go
//
// TODO: properly handle ansi escape sequences when wrapping colored strings
// into multiple lines.
func WrapWords(s string, width int) string {
	words := strings.Fields(s)

	if len(words) == 0 {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(words[0])

	spaceLeft := width - displayWidth(words[0])

	for _, word := range words[1:] {
		wordWidth := displayWidth(word)
		if wordWidth+1 > spaceLeft {
			sb.WriteRune(newline)
			sb.WriteString(word)
			spaceLeft = width - wordWidth
			continue
		}

		sb.WriteRune(space)
		sb.WriteString(word)
		spaceLeft -= 1 + wordWidth
	}

	return sb.String()
}
