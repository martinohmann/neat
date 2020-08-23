package main

import (
	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/style"
)

func main() {
	console.Print(
		"This ", style.FgRed, "line ", style.Reset, "contains ",
		style.New(style.FgBlack, style.BgYellow), "some", style.Reset,
		style.Bold, " styling examples.\n",
	)

	console.Println("{green}Something {reset,bold,underline}like this{rst} is also {bgblue,fgblack}possible{rst}.")

	blueish := style.New(style.BgRGB(0, 100, 200))

	blueish.Println("You can also style complete lines...")

	console.Print("{bgyellow,fgblack}...or combine", style.Reset, " {blue}several{rst} ", style.Fg256(200), "approaches.\n")

	style.AttributeMap["custom"] = style.New(style.FgRGB(100, 0, 100), style.BgRGB(100, 100, 0))

	console.Println("How about defining a {custom}custom named style{reset}?")

	console.Printf("%vstyleing%v also works in {bold}format strings!\n", style.FgMagenta, style.FgCyan)

	console.Println(style.BgRed, "Some", style.BgGreen, "Println", "example", style.BgYellow, style.FgBlack, "with styles")
}
