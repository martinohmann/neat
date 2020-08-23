package text

import (
	"github.com/martinohmann/neat/internal/util"
	"github.com/martinohmann/neat/measure"
	"github.com/martinohmann/neat/style"
)

// Text is a console.Renderable that produces aligned and styled text.
type Text struct {
	// Alignment controls the alignment of the text. Default is to align left.
	Alignment Alignment
	// Style is applied to the text after alignment and word wrapping.
	Style *style.Style
	// Text contains the text that should be rendered. Can contain newlines or
	// even ANSI escape sequences.
	Text string
	// WordWrap controls the word wrapping behaviour. If true, words are
	// wrapped onto multiple lines depending on the desired render width.
	WordWrap bool
}

// New creates a new Text.
func New(text string) Text {
	return Text{
		Text:      text,
		Alignment: AlignLeft,
	}
}

// Measure implements console.Renderable.
func (t Text) Measure(maxWidth int) measure.Measurement {
	width := DisplayWidth(t.maybeWordWrap(maxWidth))

	width = util.MinInt(width, maxWidth)

	return measure.NewMeasurement(width, width)
}

// Render implements console.Renderable.
func (t Text) Render(width int) string {
	text := Align(t.maybeWordWrap(width), width, t.Alignment)

	lines := SplitLines(text)
	for i, line := range lines {
		line = Truncate(line, width)

		if t.Style != nil {
			line = t.Style.Sprint(line)
		}

		lines[i] = line
	}

	return JoinLines(lines)
}

func (t Text) maybeWordWrap(width int) string {
	if t.WordWrap {
		return WrapWords(t.Text, width)
	}
	return t.Text
}
