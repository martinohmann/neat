package main

import (
	"fmt"
	"time"

	"github.com/martinohmann/neat/bar"
	"github.com/martinohmann/neat/progress"
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/text"
)

func main() {
	const numTasks = 5

	opts := []progress.Option{
		progress.WithColumns(
			progress.NewTextFuncColumn(func(task *progress.Task) string { return task.Description() }, style.New(style.FgBlue), text.AlignCenter),
			progress.NewTextFuncColumn(func(task *progress.Task) string { return task.Description() }, nil, text.AlignLeft),
			progress.NewTextFuncColumn(func(task *progress.Task) string { return task.Description() }, nil, text.AlignJustify),
			progress.TextColumn{
				Alignment: text.AlignJustify,
				Style:     style.New(style.FgYellow),
				TextFunc:  func(task *progress.Task) string { return task.Description() },
				WordWrap:  true,
			},
			progress.NewBarColumn(),
			progress.NewProgressColumn(),
			progress.NewETAColumn(),
			progress.NewTextColumn("{}", style.New(style.FgCyan), text.AlignCenter),
			progress.BarColumn{
				MaxWidth:       15,
				FinishedStyle:  bar.NewStyle('❯', style.New(style.FgGreen, style.BgBlack)),
				CompletedStyle: bar.NewStyle('❯', style.New(style.FgRed)),
			},
			progress.NewPercentageColumn(),
		),
	}

	p := progress.New(opts...)
	defer p.Stop()

	tasks := make([]*progress.Task, numTasks)

	for i := 0; i < numTasks; i++ {
		task := p.NewTask(fmt.Sprintf("task %d\nyay", i), int64((i+1)*100))

		tasks[i] = task
	}

	for i := 0; i < numTasks; i++ {
		go executeTask(tasks[i], i+1)
	}

	<-time.After(2 * time.Second)

	task := p.NewTask(fmt.Sprintf("extraextra\nlong task %d", numTasks), int64((numTasks+5)*100))

	<-time.After(2 * time.Second)

	go executeTask(task, numTasks)
	go executeTask(task, numTasks)
	go executeTask(task, numTasks)
	go executeTask(task, numTasks)

	p.Wait()

	fmt.Println("finished")
}

func executeTask(task *progress.Task, i int) {
	for !task.Finished() {
		task.Advance(int64(i))
		<-time.After(time.Duration(i*20) * time.Millisecond)
	}
}
