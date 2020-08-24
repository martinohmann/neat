package main

import (
	"fmt"
	"time"

	"github.com/martinohmann/neat/progress"
)

func main() {
	for i := 0; i < 5; i++ {
		run(i)
	}
}

func run(i int) {
	p := progress.New()
	defer p.Stop()

	task := p.NewTask(fmt.Sprintf("example %d", i), 200)

	for i := 0; i < 200; i++ {
		task.Advance(1)
		time.Sleep(2000 * time.Millisecond)
	}
}
