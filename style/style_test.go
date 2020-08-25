package style

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyle_Print(t *testing.T) {
	defer Enable()()
	assert := assert.New(t)

	assert.Equal("\x1b[1mfoo\x1b[0m", New(Bold).Sprint("foo"))
	assert.Equal("\x1b[32;1mfoo 42\x1b[0m", New(FgGreen, Bold).Sprintf("foo %d", 42))
	assert.Equal("\x1b[38;5;244mfoo\n\x1b[0m", New(Fg256(244)).Sprintln("foo"))
	assert.Equal("\x1b[38;2;64;128;255mfoo\x1b[0m", New(FgRGB(64, 128, 255)).Sprint("foo"))
	assert.Equal("\x1b[48;2;64;128;255mfoo\x1b[0m", New(BgHex(0x4080FF)).Sprint("foo"))

	var buf bytes.Buffer

	New(FgHex(0xAABBCC)).Fprint(&buf, "foo")

	assert.Equal("\x1b[38;2;170;187;204mfoo\x1b[0m", buf.String())
}
