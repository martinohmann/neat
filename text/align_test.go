package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlign(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("foo  ", Align("foo", 5, AlignLeft))
	assert.Equal("foo    \nbarbaz ", Align("foo\nbarbaz", 7, AlignLeft))
	assert.Equal("foobarbaz", Align("foobarbaz", 5, AlignLeft))
	assert.Equal("  foo", Align("foo", 5, AlignRight))
	assert.Equal("    foo\n barbaz", Align("foo\nbarbaz", 7, AlignRight))
	assert.Equal("foobarbaz", Align("foobarbaz", 5, AlignRight))
	assert.Equal(" foo ", Align("foo", 5, AlignCenter))
	assert.Equal(" foo  ", Align("foo", 6, AlignCenter))
	assert.Equal("  foo   \n barbaz ", Align("foo\nbarbaz", 8, AlignCenter))
	assert.Equal("foobarbaz", Align("foobarbaz", 5, AlignCenter))
	assert.Equal("foo bar baz\n          ", Align("foo bar baz\n", 10, AlignJustify))
	assert.Equal("foo bar baz\nqux       ", Align("foo bar baz\nqux", 10, AlignJustify))
	assert.Equal("foo bar baz", Align("foo bar baz", 10, AlignJustify))
	assert.Equal("foo bar baz", Align("foo bar baz", 11, AlignJustify))
	assert.Equal("foo bar  baz", Align("foo bar baz", 12, AlignJustify))
	assert.Equal("foo  bar  baz", Align("foo bar baz", 13, AlignJustify))
	assert.Equal("foo  bar   baz", Align("foo bar baz", 14, AlignJustify))
	assert.Equal("foo   bar   baz", Align("foo bar baz", 15, AlignJustify))
	assert.Equal("foo  barbaz\nbaz     qux\nlorem ipsum", Align("foo barbaz\nbaz qux\nlorem ipsum", 11, AlignJustify))
}
