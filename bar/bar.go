package bar

import (
	"github.com/martinohmann/neat/internal/util"
	"github.com/martinohmann/neat/measure"
)

// Bar renders a progress bar.
type Bar struct {
	RemainingStyle *Style
	CompletedStyle *Style
	FinishedStyle  *Style
	MaxWidth       int
	Completed      float64
}

// New creates a new Bar which is completed by the specified percentage.
func New(completed float64) Bar {
	return Bar{
		Completed: completed,
		MaxWidth:  -1,
	}
}

// Measure implements console.Renderable.
func (b Bar) Measure(maxWidth int) measure.Measurement {
	maximum := maxWidth
	if b.MaxWidth >= 0 {
		maximum = util.MinInt(b.MaxWidth, maximum)
	}

	return measure.NewMeasurement(4, util.MaxInt(4, maximum))
}

// Render implements console.Renderable.
func (b Bar) Render(width int) string {
	if width <= 0 {
		return ""
	}

	completedPerc := util.MinFloat64(100, util.MaxFloat64(0, b.Completed))

	remainingStyle, completedStyle, finishedStyle := b.getStyles()

	if completedPerc == 100 {
		return finishedStyle.Render(width)
	}

	completedWidth := int(float64(width) * completedPerc / 100)
	remainingWidth := width - completedWidth

	completed := completedStyle.Render(completedWidth)
	remaining := remainingStyle.Render(remainingWidth)

	return completed + remaining
}

func (b Bar) getStyles() (remaining, completed, finished *Style) {
	remaining = b.RemainingStyle
	if remaining == nil {
		remaining = DefaultRemainingStyle
	}

	completed = b.CompletedStyle
	if completed == nil {
		completed = DefaultCompletedStyle
	}

	finished = b.FinishedStyle
	if finished == nil {
		finished = DefaultFinishedStyle
	}

	return remaining, completed, finished
}
