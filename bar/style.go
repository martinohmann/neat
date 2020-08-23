package bar

import (
	"fmt"
	"strings"

	"github.com/martinohmann/neat/measure"
	"github.com/martinohmann/neat/style"
	runewidth "github.com/mattn/go-runewidth"
)

// Default progress bar styles. These are used if the corresponding styles of a
// Bar are not explicitly set.
var (
	DefaultRemainingStyle = NewStyle('─', style.New(style.FgBlack))
	DefaultCompletedStyle = NewStyle('─', style.New(style.FgRed))
	DefaultFinishedStyle  = NewStyle('─', style.New(style.FgGreen))
)

// Style is the style of a progress bar.
type Style struct {
	style  *style.Style
	symbol rune
}

// NewStyle creates a new *Style. Will panic if symbol does not have a rune
// width of 1.
func NewStyle(symbol rune, style *style.Style) *Style {
	width := runewidth.RuneWidth(symbol)
	if width != 1 {
		panic(fmt.Sprintf("NewStyle: symbol must have a rune width of 1, got %d", width))
	}

	return &Style{style, symbol}
}

// Measure implements console.Renderable.
func (s *Style) Measure(_ int) measure.Measurement {
	return measure.NewMeasurement(1, 1)
}

// Render implements console.Renderable.
func (s *Style) Render(width int) string {
	if width <= 0 {
		return ""
	}

	bar := strings.Repeat(string(s.symbol), width)

	if s.style != nil {
		bar = s.style.Sprint(bar)
	}

	return bar
}
