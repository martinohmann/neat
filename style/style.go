package style

import (
	"fmt"
	"io"
	"os"
	"strings"

	colorable "github.com/mattn/go-colorable"
	isatty "github.com/mattn/go-isatty"
)

const escape = "\x1b"

var (
	// Stdout is an io.Writer for stdout which properly handles escape
	// sequences.
	Stdout = colorable.NewColorableStdout()

	// colorsEnabled controls whether the output is colorized or not. This is
	// automatically set to true if the terminal is "dumb" or if stdout is not
	// a TTY. Can be overridden to forcefully enable of disable colored output.
	colorsEnabled = !(os.Getenv("TERM") == "dumb" || !isTerminal(os.Stdout.Fd()))
)

// Enabled returns true if colors are enabled.
func Enabled() bool { return colorsEnabled }

// Enable enables output coloring. The returned func can be used in combination
// with defer to restore the previous state.
func Enable() func() { return enable(true) }

// Disable disables output coloring. The returned func can be used in
// combination with defer to restore the previous state.
func Disable() func() { return enable(false) }

func enable(enabled bool) func() {
	oldColorsEnabled := colorsEnabled
	colorsEnabled = enabled

	if enabled != oldColorsEnabled {
		// Clear sequence cache if colors were enabled or disabled.
		sequenceCache.Range(func(key, value interface{}) bool {
			sequenceCache.Delete(key)
			return true
		})
	}

	return func() { enable(oldColorsEnabled) }
}

// Style can style and color text.
type Style struct {
	attrs []Attribute
}

// New creates a new *Style from given attributes.
func New(attrs ...Attribute) *Style {
	return newStyle(attrs)
}

func newStyle(attrs []Attribute) *Style {
	c := &Style{make([]Attribute, 0, len(attrs))}
	return c.add(attrs)
}

// Fg256 creates a foreground 256color attribute.
func Fg256(color uint8) *Style {
	return New(FgColor, colorMode256, SimpleAttribute(color))
}

// Bg256 creates a background 256color attribute.
func Bg256(color uint8) *Style {
	return New(BgColor, colorMode256, SimpleAttribute(color))
}

// FgRGB creates a foreground RGB color attribute.
func FgRGB(r, g, b uint8) *Style {
	return New(FgColor, colorModeRGB, SimpleAttribute(r), SimpleAttribute(g), SimpleAttribute(b))
}

// BgRGB creates a background RGB color attribute.
func BgRGB(r, g, b uint8) *Style {
	return New(BgColor, colorModeRGB, SimpleAttribute(r), SimpleAttribute(g), SimpleAttribute(b))
}

// FgHex creates a foreground RGB color attribute from a hex value.
func FgHex(v uint32) *Style {
	return FgRGB(toRGB(v))
}

// BgHex creates a foreground RGB color attribute from a hex value.
func BgHex(v uint32) *Style {
	return BgRGB(toRGB(v))
}

func toRGB(v uint32) (r uint8, g uint8, b uint8) {
	r, g, b = uint8(v>>16), uint8(v>>8), uint8(v)
	return
}

// NewWith creates a new *Style from s with additional attributes. Style s is
// not altered.
func (s *Style) NewWith(attrs ...Attribute) *Style {
	n := s.Copy()
	return n.add(attrs)
}

func (s *Style) Copy() *Style {
	return newStyle(s.attrs)
}

// Add adds attributes to an existing style.
func (s *Style) Add(attrs ...Attribute) *Style {
	return s.add(attrs)
}

func (s *Style) add(attrs []Attribute) *Style {
	for _, attr := range attrs {
		if style, ok := attr.(*Style); ok {
			s.attrs = append(s.attrs, style.attrs...)
		} else {
			s.attrs = append(s.attrs, attr)
		}
	}

	return s
}

func (s *Style) sequence() string {
	var sb strings.Builder

	for i, attr := range s.attrs {
		if i != 0 {
			sb.WriteRune(';')
		}

		sb.WriteString(attr.sequence())
	}

	return sb.String()
}

// Print formats using the default formats for its operands and writes to
// standard output. Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (s *Style) Print(args ...interface{}) (n int, err error) {
	return s.Fprint(Stdout, args...)
}

// Println formats using the default formats for its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended. It returns the number of bytes written and any write error
// encountered.
func (s *Style) Println(args ...interface{}) (n int, err error) {
	return s.Fprintln(Stdout, args...)
}

// Printf formats according to a format specifier and writes to standard
// output. It returns the number of bytes written and any write error
// encountered.
func (s *Style) Printf(format string, args ...interface{}) (n int, err error) {
	return s.Fprintf(Stdout, format, args...)
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string. It returns the
// number of bytes written and any write error encountered.
func (s *Style) Fprint(w io.Writer, args ...interface{}) (n int, err error) {
	return s.wrapWriter(w, func() (int, error) {
		return fmt.Fprint(w, args...)
	})
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended. It
// returns the number of bytes written and any write error encountered.
func (s *Style) Fprintln(w io.Writer, args ...interface{}) (n int, err error) {
	return s.wrapWriter(w, func() (int, error) {
		return fmt.Fprintln(w, args...)
	})
}

// Fprintf formats according to a format specifier and writes to w. It returns
// the number of bytes written and any write error encountered.
func (s *Style) Fprintf(w io.Writer, format string, args ...interface{}) (n int, err error) {
	return s.wrapWriter(w, func() (int, error) {
		return fmt.Fprintf(w, format, args...)
	})
}

// Sprint formats using the default formats for its operands and returns the
// resulting string. Spaces are added between operands when neither is a
// string.
func (s *Style) Sprint(args ...interface{}) string {
	return s.wrapString(func() string {
		return fmt.Sprint(args...)
	})
}

// Sprintln formats using the default formats for its operands and returns the
// resulting string. Spaces are always added between operands and a newline is
// appended.
func (s *Style) Sprintln(args ...interface{}) string {
	return s.wrapString(func() string {
		return fmt.Sprintln(args...)
	})
}

// Sprintf formats according to a format specifier and returns the resulting
// string.
func (s *Style) Sprintf(format string, args ...interface{}) string {
	return s.wrapString(func() string {
		return fmt.Sprintf(format, args...)
	})
}

func (s *Style) wrapWriter(w io.Writer, fn func() (int, error)) (n int, err error) {
	n, err = EscapeWriter(w, s)
	if err != nil {
		return
	}

	nn, err := fn()
	n += nn
	if err != nil {
		return
	}

	nn, err = ResetWriter(w)
	n += nn
	return
}

func (s *Style) wrapString(fn func() string) string {
	return EscapeString(s) + fn() + ResetString()
}

// EscapeString creates the escape sequence for given attribute and returns it.
// If coloring is disabled this returns an empty string.
func EscapeString(attr Attribute) string {
	if !colorsEnabled {
		return ""
	}

	return fmt.Sprintf("%s[%sm", escape, attr.sequence())
}

// EscapeWriter creates the escape sequence for given attribute and writes to
// w. It returns the number of bytes written and any write error encountered.
// If coloring is disabled this is a no-op.
func EscapeWriter(w io.Writer, attr Attribute) (n int, err error) {
	if !colorsEnabled {
		return
	}

	return fmt.Fprintf(w, "%s[%sm", escape, attr.sequence())
}

// ResetString creates the escape sequence for resetting all style attributes
// and returns it.
// If coloring is disabled this returns an empty string.
func ResetString() string {
	return EscapeString(Reset)
}

// ResetWriter creates the escape sequence for resetting all style attributes
// and writes to w. It returns the number of bytes written and any write error
// encountered. If coloring is disabled this is a no-op.
func ResetWriter(w io.Writer) (n int, err error) {
	return EscapeWriter(w, Reset)
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}
