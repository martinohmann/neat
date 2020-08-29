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
	console.Printf("1. {bold}lorem ipsum with colored borders and margin\n\n")

	opts := []table.Option{
		table.WithPadding(1),
		table.WithMargin(2),
		table.WithColumnAlignment(text.AlignLeft, text.AlignJustify, text.AlignRight),
		table.WithAlignment(text.AlignRight),
		table.WithWordWrap(true),
		table.WithBorderMask(table.BorderAll),
		table.WithBorderStyle(style.New(style.Bold, style.FgBlack)),
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
	t.AddRow(lorem, lorem, lorem)
	t.AddRow(left.Text[:100], center, right)

	t.Render()

	console.Printf("\n{bold}2. table with only column and bottom borders and custom border rune\n\n")

	t = table.New(
		os.Stdout,
		table.WithBorderMask(table.BorderColumn|table.BorderBottom),
		table.WithBorderStyle(style.New(style.FgRGB(200, 100, 0))),
		table.WithBorderRunes(table.BorderRunes{
			table.BorderRuneVertical:           '║',
			table.BorderRuneIntersectionBottom: '╨',
		}),
	)

	t.AddRow(1, 2, 3)
	t.AddRow("one", "two", "three")

	t.Render()

	console.Printf("\n{bold}3. simple table, no borders\n\n")

	t = table.New(os.Stdout)

	t.AddRow("FOO", "BAR", "BAZ")
	t.AddRow("ten", "eleven", "twelve")

	t.Render()
}
