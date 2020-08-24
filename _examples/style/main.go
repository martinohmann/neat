package main

import (
	"os"

	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/style"
	colorable "github.com/mattn/go-colorable"
)

func main() {
	red := style.New(style.BgRed)

	bold256 := style.New(style.Bold, style.Bg256(128))

	console.Printf("%sHello World\n", red)
	console.Printf("Hello %sWorld!\n", style.BgRGB(255, 128, 0))
	console.Println("{yellow,bold,underline}Hello{reset,green} World!")
	console.Printf("Hello %sWorld!\n", style.Bg256(100))
	console.Printf("%vHello %sWorld!\n", bold256, style.Reset)

	w := colorable.NewNonColorable(os.Stdout)

	console.Fprintln(w, "{red}bar")

	printer := console.NewPrinter(style.Stdout)

	style.AttributeMap["nobold"] = style.Normal

	printer.Println("{fgcyan,bgred}foo{bold,black}bar{nobold,yellow,bggreen}baz")

	console.Printf("%shex\n", style.BgHex(0x5588bb))
	console.Printf("%srgb\n", style.BgRGB(toRGB(0x5588bb)))

	console.Print(style.New(style.BgRed, style.FgBlack), "red", style.BgBlue, "blue", "{bggreen}green", "alsogreen\n")

	style.New(style.FgBlack, style.Bold, style.BgRGB(0, 255, 0)).Fprintln(style.Stdout, "this should be green")
}

func toRGB(v uint32) (uint8, uint8, uint8) {
	return uint8((v >> 16) & 0xff), uint8((v >> 8) & 0xff), uint8(v & 0xff)
}
