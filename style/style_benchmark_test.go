package style

import (
	"io"
	"io/ioutil"
	"testing"
)

func init() {
	Enable()
}

func BenchmarkStyleStringNoColor(b *testing.B) {
	benchmarkStyleString(b, "some string")
}

func BenchmarkStyleStringNoMatch(b *testing.B) {
	benchmarkStyleString(b, "some {}string")
}

func BenchmarkStyleStringShort(b *testing.B) {
	benchmarkStyleString(b, "{red}string")
}

func BenchmarkStyleStringMid(b *testing.B) {
	benchmarkStyleString(b, "{red}string{bold}string")
}

func BenchmarkStyleStringMidInvalid(b *testing.B) {
	benchmarkStyleString(b, "{red}string{inval}string")
}

func BenchmarkStyleStringLong(b *testing.B) {
	benchmarkStyleString(b, "{red}string{bold}string{green,bgblue}string")
}

func BenchmarkStyleStringLongInvalid(b *testing.B) {
	benchmarkStyleString(b, "{red}string{bold}string{inval,bgblue}string")
}

func benchmarkStyleString(b *testing.B, s string) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		StyleString(s)
	}
}

func BenchmarkAttributesSequence(b *testing.B) {
	a := New(FgRed, Bold, BgRGB(255, 255, 255))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a.sequence()
	}
}

func BenchmarkColorFprint(b *testing.B) {
	c := New(FgRed, Bold, BgRGB(0, 255, 0))

	benchmarkFprint(b, c.Fprint, "{green}{bgblue}blue{red}", FgRed, BgGreen, "string", FgBlack, 1, FgYellow)
}

func BenchmarkColorFprintln(b *testing.B) {
	c := New(FgRed, Bold, BgRGB(0, 255, 0))

	benchmarkFprint(b, c.Fprintln, "{green}{bgblue}blue{red}", FgRed, BgGreen, "string", FgBlack, 1, FgYellow)
}

func benchmarkFprint(b *testing.B, fprint func(io.Writer, ...interface{}) (int, error), args ...interface{}) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fprint(ioutil.Discard, args...)
	}
}

func BenchmarkColorFprintf(b *testing.B) {
	c := New(FgRed, Bold, BgRGB(0, 255, 0))

	benchmarkFprintf(b, c.Fprintf, "{green}%s{bgblue}blue{red}%v%d", "{bold}string{reset}", FgBlack, 1)
}

func benchmarkFprintf(b *testing.B, fprintf func(io.Writer, string, ...interface{}) (int, error), format string, args ...interface{}) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fprintf(ioutil.Discard, format, args...)
	}
}
