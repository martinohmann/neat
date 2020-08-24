package progress

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/table"
)

// Progress manages and display task progress information.
type Progress struct {
	out console.FileWriter

	columns      []Column
	tableOptions []table.Option

	tasks []*Task

	addCh    chan *Task
	notifyCh chan struct{}
	stopCh   chan struct{}
	waitCh   chan struct{}

	stopped int32

	displayHeight int
}

// New creates a new *Progress with given options.
func New(opts ...Option) *Progress {
	p := &Progress{
		addCh:    make(chan *Task),
		notifyCh: make(chan struct{}),
		stopCh:   make(chan struct{}),
		waitCh:   make(chan struct{}),
	}

	for _, option := range opts {
		option(p)
	}

	if p.out == nil {
		p.out = os.Stdout
	}

	if p.columns == nil {
		p.columns = DefaultColumns
	}

	go p.run()

	return p
}

// NewTask creates a new *Task that will be added to the progress using the
// provided description and total. If the progress is already stopped this will
// just create a task but does not add it to the progress.
func (p *Progress) NewTask(desc string, total int64) *Task {
	task := newTask(p, desc, total)

	// Only accept new tasks if not stopped. This avoid deadlocks when the
	// progress is stopped but new tasks are still created.
	if !p.Stopped() {
		p.addCh <- task
	}

	return task
}

func (p *Progress) notify() {
	select {
	case p.notifyCh <- struct{}{}:
	default:
		// Avoid blocked tasks because Stop was called.
	}
}

// Wait waits until all tasks are finished. Must be called after all desired
// tasks are created. Otherwise the progress is stopped the first time all
// currently created tasks are finished.
func (p *Progress) Wait() {
	<-p.waitCh
}

// Stop stops the progress. Should be called using defer to ensure that the
// terminal cursor is reset after the progress is finished or the program was
// interrupted. It is safe to call Stop multiple times, subsequent calls are
// no-op.
func (p *Progress) Stop() {
	if !atomic.CompareAndSwapInt32(&p.stopped, 0, 1) {
		return
	}

	close(p.stopCh)
}

// Stopped returns true if the progress is stopped.
func (p *Progress) Stopped() bool {
	return atomic.LoadInt32(&p.stopped) == 1
}

func (p *Progress) run() {
	cursor := &terminal.Cursor{Out: p.out}

	defer func() {
		// Reset the cursor and unblock all waiters.
		cursor.HorizontalAbsolute(0)
		cursor.Show()
		close(p.waitCh)
	}()

	cursor.Hide()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var finished bool

	for {
		select {
		case task := <-p.addCh:
			p.tasks = append(p.tasks, task)
		case <-p.notifyCh:
			finished = p.update(cursor)
		case <-ticker.C:
			finished = p.update(cursor)
		case <-p.stopCh:
			return
		}

		if finished {
			return
		}
	}
}

// update renders the progress for all started tasks and repositions the
// cursor. Returns true if all tasks are finished.
func (p *Progress) update(cursor *terminal.Cursor) bool {
	if len(p.tasks) == 0 {
		return false
	}

	cursor.HorizontalAbsolute(0)
	if p.displayHeight > 0 {
		cursor.Up(p.displayHeight)
	}

	table := table.New(p.out, p.tableOptions...)

	for _, task := range p.tasks {
		if !task.Started() {
			continue
		}

		table.AddRow(p.tableColumns(task)...)
	}

	tableHeight, err := table.Render()
	if err != nil {
		panic(fmt.Errorf("failed to render task progress: %v", err))
	}

	p.displayHeight = tableHeight

	return p.finished()
}

// finished returns true if all tasks for this progress are finished.
func (p *Progress) finished() bool {
	for _, task := range p.tasks {
		if !task.Finished() {
			return false
		}
	}

	return true
}

// tableColumns renders progress columns for given tasks into a slice of
// interface{} that can be added as a table row.
func (p *Progress) tableColumns(task *Task) []interface{} {
	cols := make([]interface{}, len(p.columns))

	for i, c := range p.columns {
		cols[i] = c.Render(task)
	}

	return cols
}
