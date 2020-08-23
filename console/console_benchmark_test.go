package console

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/martinohmann/neat/style"
)

func init() {
	style.Enable()
}

func BenchmarkFmtFprint(b *testing.B) {
	benchmarkFprint(b, fmt.Fprint, "{green}{bgblue}blue{red}", style.FgRed, style.BgGreen, "string", style.FgBlack, 1, style.FgYellow)
}

func BenchmarkFprint(b *testing.B) {
	benchmarkFprint(b, Fprint, "{green}{bgblue}blue{red}", style.FgRed, style.BgGreen, "string", style.FgBlack, 1, style.FgYellow)
}

func BenchmarkFmtFprintln(b *testing.B) {
	benchmarkFprint(b, fmt.Fprintln, "{green}{bgblue}blue{red}", style.FgRed, style.BgGreen, "string", style.FgBlack, 1, style.FgYellow)
}

func BenchmarkFprintln(b *testing.B) {
	benchmarkFprint(b, Fprintln, "{green}{bgblue}blue{red}", style.FgRed, style.BgGreen, "string", style.FgBlack, 1, style.FgYellow)
}

func benchmarkFprint(b *testing.B, fprint func(io.Writer, ...interface{}) (int, error), args ...interface{}) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fprint(ioutil.Discard, args...)
	}
}

func BenchmarkFmtFprintf(b *testing.B) {
	benchmarkFprintf(b, fmt.Fprintf, "{green}%s{bgblue}blue{red}%v%d", "{bold}string{reset}", style.FgBlack, 1)
}

func BenchmarkFprintf(b *testing.B) {
	benchmarkFprintf(b, Fprintf, "{green}%s{bgblue}blue{red}%v%d", "{bold}string{reset}", style.FgBlack, 1)
}

func benchmarkFprintf(b *testing.B, fprintf func(io.Writer, string, ...interface{}) (int, error), format string, args ...interface{}) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fprintf(ioutil.Discard, format, args...)
	}
}
