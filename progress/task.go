package progress

import (
	"sync"
	"time"
)

// timeNow be overwritten with a fake now provider in tests.
var timeNow = time.Now

// Task is a container for the state of a task's current progress.
type Task struct {
	mu        sync.Mutex
	progress  *Progress
	desc      string
	startTime time.Time
	total     int64
	completed int64
}

// newTask creates a new *Task with description and total. Will panic if total
// is < 0.
func newTask(progress *Progress, desc string, total int64) *Task {
	if total < 0 {
		panic("total must be greater than or equal to 0")
	}

	return &Task{
		desc:     desc,
		total:    total,
		progress: progress,
	}
}

// Advance advances the progress of given task by step.
func (t *Task) Advance(step int64) {
	t.mu.Lock()

	if t.startTime.IsZero() {
		t.startTime = timeNow()
	}

	if t.completed+step <= t.total {
		t.completed += step
	} else {
		t.completed = t.total
	}

	t.mu.Unlock()

	t.progress.notify()
}

// Start starts a task by setting its start time.
func (t *Task) Start() {
	t.Advance(0)
}

// Description returns the tasks description.
func (t *Task) Description() string {
	return t.desc
}

// Completed returns the number of completed parts of the task. This is always
// less than or equal to the total.
func (t *Task) Completed() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.completed
}

// Total returns the total number of parts of the task.
func (t *Task) Total() int64 {
	return t.total
}

// Elapsed returns the elapsed time for the task.
func (t *Task) Elapsed() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return timeNow().Sub(t.startTime)
}

// Estimated returns the estimated duration until the task is finished based on
// the rate of change.
func (t *Task) Estimated() time.Duration {
	t.mu.Lock()
	completed := t.completed
	startTime := t.startTime
	t.mu.Unlock()

	now := timeNow()
	ratio := float64(completed) / float64(t.total)
	elapsed := now.Sub(startTime)
	total := time.Duration(float64(elapsed) / ratio)
	estimated := startTime.Add(total)

	return estimated.Sub(now)
}

// PercentCompleted returns the completed percentage of the task as a value 0
// <= percentage <= 100.
func (t *Task) PercentCompleted() float64 {
	return float64(t.Completed()) / float64(t.total) * 100
}

// Started returns true if the task was already started. A task is started if
// either the Start or Advance method was called at least once.
func (t *Task) Started() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.startTime.IsZero()
}

// Finished returns true if the task is finished, that is, the total is
// reached.
func (t *Task) Finished() bool {
	return t.Completed() >= t.total
}
