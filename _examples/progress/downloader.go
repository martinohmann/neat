package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/martinohmann/neat/progress"
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/text"
)

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		return
	}

	p := progress.New(progress.WithColumns(
		progress.NewDescriptionColumn(),
		progress.NewBarColumn(),
		progress.TextColumn{
			Alignment: text.AlignRight,
			Style:     style.New(style.FgCyan),
			TextFunc: func(task *progress.Task) string {
				completed := uint64(task.Completed())
				total := uint64(task.Total())
				return fmt.Sprintf("%s/%s", humanize.Bytes(completed), humanize.Bytes(total))
			},
		},
		progress.TextColumn{Text: "•", Style: style.New(style.FgHex(0x5e81ac))},
		progress.NewPercentageColumn(),
		progress.TextColumn{Text: "•", Style: style.New(style.FgYellow)},
		progress.NewETAColumn(),
	))
	defer p.Stop()

	errors := make([]error, 0, len(urls))
	errCh := make(chan error, len(urls))

	for _, url := range urls {
		go downloadFile(p, errCh, url)
	}

	for err := range errCh {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		fmt.Printf("download errors: %v\n", errors)
	}
}

func downloadFile(p *progress.Progress, errCh chan<- error, url string) {
	headResp, err := http.Head(url)
	if err != nil {
		errCh <- err
		return
	}
	defer headResp.Body.Close()

	size, err := strconv.ParseInt(headResp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errCh <- err
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	task := p.NewTask(filepath.Base(url), size)

	reader := task.ReadCounter(resp.Body)

	task.Start()

	_, err = io.Copy(ioutil.Discard, reader)

	errCh <- err
}
