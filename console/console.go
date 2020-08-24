package console

import (
	"io"

	"github.com/martinohmann/neat/measure"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	defaultWidth  = 80
	defaultHeight = 25
)

// Renderable is an object that can be rendered to a console.
type Renderable interface {
	// Measure returns a width measurement which indicates the minimum and
	// maximum widths required to correctly render the Renderable. The passed
	// in maxWidth value indicates the maximum usable width.
	Measure(maxWidth int) measure.Measurement

	// Render renders the Renderable using the given width. If width is not
	// within the bounds returned by Measure, the rendered result may be
	// undefined.
	Render(width int) string
}

// FileWriter is an io.Writer which also provides access to the file
// descriptor. The interface is satisfied by *os.File.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// TerminalSize returns the terminal width and height of given FileWriter.
func TerminalSize(fw FileWriter) (w, h int) {
	width, height, err := terminal.GetSize(int(fw.Fd()))
	if err != nil {
		return defaultWidth, defaultHeight
	}

	return width, height
}

// TerminalWidth returns the terminal width of given FileWriter.
func TerminalWidth(fw FileWriter) int {
	width, _ := TerminalSize(fw)
	return width
}

// TerminalHeight returns the terminal height of given FileWriter.
func TerminalHeight(fw FileWriter) int {
	_, height := TerminalSize(fw)
	return height
}
