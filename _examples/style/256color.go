package main

import (
	"fmt"

	"github.com/martinohmann/neat/style"
)

func main() {
	fmt.Printf("%-45s 256 Color(16 bit) Table %-35s\n", " ", " ")
	// 0 - 16
	fmt.Printf("%-22sStandard Color %-42sExtended Color \n", " ", " ")

	var fg uint8 = 255
	for i := 0; i < 8; i++ {
		if i > 3 {
			fg = 0
		}
		style.New(style.Fg256(fg), style.Bg256(uint8(i))).Printf("   %-4d", i)
	}

	fmt.Print("    ")

	fg = 255
	for i := 8; i < 16; i++ {
		if i > 11 {
			fg = 0
		}
		style.New(style.Fg256(fg), style.Bg256(uint8(i))).Printf("   %-4d", i)
	}

	fg = 255
	fmt.Printf("\n%-50s216 Color\n", " ")
	for i := 16; i < 232; i++ {
		v := i - 16

		if v != 0 {
			if v%18 == 0 {
				fg = 0
				fmt.Println()
			}

			if v%36 == 0 {
				fg = 255
				fmt.Println()
			}
		}

		style.New(style.Fg256(fg), style.Bg256(uint8(i))).Printf("  %-4d", i)
	}

	fmt.Printf("\n%-50s24th Order Grayscale Color\n", " ")
	for i := 232; i < 256; i++ {
		v := i - 232
		if v < 12 {
			fg = 255
		} else {
			fg = 0
		}

		style.New(style.Fg256(fg), style.Bg256(uint8(i))).Printf("  %-4d", i)
	}
	fmt.Println()
}
