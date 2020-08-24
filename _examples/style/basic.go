package main

import (
	"fmt"
	"sort"

	"github.com/martinohmann/neat/console"
	"github.com/martinohmann/neat/style"
)

func main() {
	names := make([]string, 0, len(style.AttributeMap))
	maxlen := 0

	for name := range style.AttributeMap {
		names = append(names, name)
		if len(name) > maxlen {
			maxlen = len(name)
		}
	}

	sort.Strings(names)

	for i, name := range names {
		if i%8 == 0 {
			fmt.Println()
		}

		attr := style.AttributeMap[name]

		// Other options:
		//   style.New(attr).Printf("%-*s", maxlen, name)
		//   console.Print(attr, fmt.Sprintf("%-*s", maxlen, name))
		console.Printf("%s%-*s", attr, maxlen, name)

		// Reset and pad.
		console.Print(style.Reset, " ")
	}
}
