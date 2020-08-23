package console

import (
	"fmt"
	"io"

	"github.com/martinohmann/neat/style"
)

// Print formats using the default formats for its operands and writes to
// standard output. Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Print(args ...interface{}) (n int, err error) {
	return Fprint(style.Stdout, args...)
}

// Println formats using the default formats for its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended. It returns the number of bytes written and any write error
// encountered.
func Println(args ...interface{}) (n int, err error) {
	return Fprintln(style.Stdout, args...)
}

// Printf formats according to a format specifier and writes to standard
// output. It returns the number of bytes written and any write error
// encountered.
func Printf(format string, args ...interface{}) (n int, err error) {
	return Fprintf(style.Stdout, format, args...)
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string. It returns the
// number of bytes written and any write error encountered.
func Fprint(w io.Writer, args ...interface{}) (n int, err error) {
	return wrapWriter(w, func() (int, error) {
		return fmt.Fprint(w, styleArgs(args)...)
	})
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended. It
// returns the number of bytes written and any write error encountered.
func Fprintln(w io.Writer, args ...interface{}) (n int, err error) {
	return wrapWriter(w, func() (int, error) {
		return fmt.Fprintln(w, styleArgs(args)...)
	})
}

// Fprintf formats according to a format specifier and writes to w. It returns
// the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, args ...interface{}) (n int, err error) {
	return wrapWriter(w, func() (int, error) {
		return fmt.Fprintf(w, style.StyleString(format), styleArgs(args)...)
	})
}

// Sprint formats using the default formats for its operands and returns the
// resulting string. Spaces are added between operands when neither is a
// string.
func Sprint(args ...interface{}) string {
	return wrapString(func() string {
		return fmt.Sprint(styleArgs(args)...)
	})
}

// Sprintln formats using the default formats for its operands and returns the
// resulting string. Spaces are always added between operands and a newline is
// appended.
func Sprintln(args ...interface{}) string {
	return wrapString(func() string {
		return fmt.Sprintln(styleArgs(args)...)
	})
}

// Sprintf formats according to a format specifier and returns the resulting
// string.
func Sprintf(format string, args ...interface{}) string {
	return wrapString(func() string {
		return fmt.Sprintf(style.StyleString(format), styleArgs(args)...)
	})
}

func wrapWriter(w io.Writer, fn func() (int, error)) (n int, err error) {
	n, err = fn()
	if err != nil {
		return
	}

	nn, err := style.ResetWriter(w)
	n += nn
	return
}

func wrapString(fn func() string) string {
	return fn() + style.ResetString()
}

func styleArgs(args []interface{}) []interface{} {
	for i, arg := range args {
		switch v := arg.(type) {
		case style.Attribute:
			args[i] = style.EscapeString(v)
		case string:
			args[i] = style.StyleString(v)
		}
	}

	return args
}

// Printer wraps an io.Writer for writing colorful strings to it.
type Printer struct {
	io.Writer
}

// NewPrinter returns a new *Printer which writes to w.
func NewPrinter(w io.Writer) *Printer {
	return &Printer{w}
}

// Print formats using the default formats for its operands and writes to the
// underlying io.Writer. Spaces are added between operands when neither is a
// string. It returns the number of bytes written and any write error
// encountered.
func (p *Printer) Print(args ...interface{}) (n int, err error) {
	return Fprint(p, args...)
}

// Println formats using the default formats for its operands and writes to to
// the underlying io.Writer. Spaces are always added between operands and a
// newline is appended. It returns the number of bytes written and any write
// error encountered.
func (p *Printer) Println(args ...interface{}) (n int, err error) {
	return Fprintln(p, args...)
}

// Printf formats according to a format specifier and writes to the underlying
// io.Writer. It returns the number of bytes written and any write error
// encountered.
func (p *Printer) Printf(format string, args ...interface{}) (n int, err error) {
	return Fprintf(p, format, args...)
}
