package progress

import (
	"fmt"

	"github.com/martinohmann/neat/bar"
	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/internal/util"
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/text"
)

var (
	// DefaultColumns contains the columns that a progress will render by default
	// for each task. Can be overridden per progress or globally by changing
	// DefaultColumns.
	DefaultColumns = []Column{
		NewDescriptionColumn(),
		NewBarColumn(),
		NewProgressColumn(),
		NewPercentageColumn(),
	}
)

// Column is a column displayed for the progress of a given task.
type Column interface {
	// Render renders task into a console.Renderable.
	Render(task *Task) console.Renderable
}

// ColumnFunc is a func that satisfies the Column interface.
type ColumnFunc func(*Task) console.Renderable

// Render implements Column.
func (f ColumnFunc) Render(task *Task) console.Renderable {
	return f(task)
}

type TextFunc func(*Task) string

type TextColumn struct {
	Alignment text.Alignment
	Style     *style.Style
	TextFunc  TextFunc
	Text      string
	WordWrap  bool
}

func (c TextColumn) Render(task *Task) console.Renderable {
	return text.Text{
		Alignment: c.Alignment,
		Style:     c.Style,
		Text:      c.text(task),
		WordWrap:  c.WordWrap,
	}
}

func (c TextColumn) text(task *Task) string {
	if c.TextFunc != nil {
		return c.TextFunc(task)
	}

	return c.Text
}

func NewTextColumn(s string, style *style.Style, alignment text.Alignment) TextColumn {
	return TextColumn{
		Alignment: alignment,
		Style:     style,
		Text:      s,
	}
}

func NewTextFuncColumn(fn TextFunc, style *style.Style, alignment text.Alignment) TextColumn {
	return TextColumn{
		Alignment: alignment,
		Style:     style,
		TextFunc:  fn,
	}
}

func NewDescriptionColumn() TextColumn {
	return TextColumn{
		Style:     style.New(style.Bold),
		Alignment: text.AlignRight,
		TextFunc: func(task *Task) string {
			return task.Description()
		},
	}
}

func NewProgressColumn() TextColumn {
	return TextColumn{
		Style:     style.New(style.FgCyan),
		Alignment: text.AlignRight,
		TextFunc: func(task *Task) string {
			digits := util.CountDigitsInt64(task.Total())
			return fmt.Sprintf("%*d/%d", digits, task.Completed(), task.Total())
		},
	}
}

func NewPercentageColumn() TextColumn {
	return TextColumn{
		Alignment: text.AlignRight,
		TextFunc: func(task *Task) string {
			return fmt.Sprintf("%3.0f%%", task.PercentCompleted())
		},
	}
}

func NewETAColumn() TextColumn {
	return TextColumn{
		Style:     style.New(style.FgGreen),
		Alignment: text.AlignRight,
		TextFunc: func(task *Task) string {
			return fmt.Sprintf("%s ETA", formatDuration(task.Estimated()))
		},
	}
}

// BarColumn is a column that will render a progress bar.
type BarColumn struct {
	RemainingStyle *bar.Style
	CompletedStyle *bar.Style
	FinishedStyle  *bar.Style
	MaxWidth       int
}

// NewBarColumn creates a new BarColumn which tried to occupy as much
// horizontal space as possible with respect to over columns.
func NewBarColumn() BarColumn {
	return BarColumn{
		MaxWidth: -1,
	}
}

// Render implements Column.
func (c BarColumn) Render(task *Task) console.Renderable {
	return bar.Bar{
		RemainingStyle: c.RemainingStyle,
		CompletedStyle: c.CompletedStyle,
		FinishedStyle:  c.FinishedStyle,
		MaxWidth:       c.MaxWidth,
		Completed:      task.PercentCompleted(),
	}
}
