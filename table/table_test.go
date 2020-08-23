package table

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTable_Render(t *testing.T) {
	assert := assert.New(t)

	table := New(40)
	table.AddRow("foo\nbarbaz", "qux")

	var buf bytes.Buffer

	nlines, err := table.Render(&buf)
	assert.NoError(err)
	assert.Equal(2, nlines)
	assert.Equal("foo    qux\nbarbaz    \n", buf.String())
}

type Suite struct {
	suite.Suite
}

func (s *Suite) testTableRender(table *Table, expected string) {
	var buf bytes.Buffer

	// remove leading newline char
	expected = strings.TrimLeft(expected, "\n")
	expected = strings.ReplaceAll(expected, ".", " ")

	_, err := table.Render(&buf)
	s.NoError(err)
	s.Equal(expected, buf.String())
}

func (s *Suite) TestTable_Render() {
	s.testTableRender(
		New(40).
			AddRow("foo\nbarbaz", "qux"),
		`
foo....qux
barbaz....
`,
	)
	s.testTableRender(
		New(40).
			AddRow("foo", "bar"),
		`
foo.bar
`,
	)
	s.testTableRender(
		New(10).
			AddRow("foo", "bar", "baz"),
		`
f….b….b…
`,
	)
	s.testTableRender(
		New(30).
			AddRow("foo", "bar", "baz").
			AddRow("foofoo", "barbar", "bazbaz").
			AddRow(1, 2, "foo"),
		`
foo....bar....baz...
foofoo.barbar.bazbaz
1......2......foo...
`,
	)
}

func TestMain(t *testing.T) {
	suite.Run(t, new(Suite))
}
