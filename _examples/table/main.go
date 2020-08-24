package main

import (
	"os"

	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/table"
	"github.com/martinohmann/neat/text"
)

const lorem = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."

func main() {
	maxWidth := console.TerminalWidth(os.Stdout)

	t := table.New(
		maxWidth,
		table.WithPadding(4),
		table.WithAlignment(text.AlignLeft, text.AlignJustify, text.AlignRight),
	)

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

	t.Render(os.Stdout)
}
