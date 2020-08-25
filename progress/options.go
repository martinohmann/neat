package progress

import (
	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/table"
)

// Option is a func for configuring a *Progress.
type Option func(p *Progress)

// WithOutput sets the output for the progress. If omitted, os.Stdout will be
// used.
func WithOutput(out console.FileWriter) Option {
	return func(p *Progress) {
		p.out = out
	}
}

// WithColumns sets the columns that the progress should render for each task.
// If omitted, DefaultColumns will be used.
func WithColumns(columns ...Column) Option {
	return func(p *Progress) {
		p.columns = columns
	}
}

// WithTableOptions sets the options for the underlying table used for
// rendering the progress of all tasks.
func WithTableOptions(opts ...table.Option) Option {
	return func(p *Progress) {
		p.tableOptions = opts
	}
}
