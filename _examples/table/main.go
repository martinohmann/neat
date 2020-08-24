package main

import (
	"fmt"
	"os"

	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/table"
	"github.com/martinohmann/neat/text"
)

const lorem = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."

func main() {
	opts := []table.Option{
		table.WithPadding(4),
		table.WithAlignment(text.AlignLeft, text.AlignJustify, text.AlignRight),
	}

	t := table.New(os.Stdout, opts...)

	bold := style.New(style.Bold)

	_ = bold

	left := text.Text{
		Text:     lorem,
		WordWrap: true,
	}

	center := text.Text{
		Text:     lorem,
		WordWrap: true,
	}

	right := text.Text{
		Text:     lorem,
		WordWrap: true,
	}

	t.AddRow(bold.Sprint("left aligned"), bold.Sprint("justify + wordwrap"), bold.Sprint("right aligned"))
	t.AddRow(left, center, right)

	n, err := t.Render()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n%d lines rendered\n\n", n)

	t = table.New(os.Stdout)

	t.AddRow(1, 2, 3)
	t.AddRow("one", "two", "three")

	n, err = t.Render()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n%d lines rendered\n", n)
}
