package table

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/martinohmann/neat/bar"
	"github.com/martinohmann/neat/style"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTable_Render_Reset(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	table := New(&buf)
	table.AddRow("foo\nbarbaz", "qux").AddRow("one", "two")

	nlines, err := table.Render()
	assert.NoError(err)
	assert.Equal(3, nlines)
	assert.Equal("foo    qux\nbarbaz    \none    two\n", buf.String())

	buf.Reset()
	table.Reset()

	table.AddRow("foo", "bar", "baz")

	nlines, err = table.Render()
	assert.NoError(err)
	assert.Equal(1, nlines)
	assert.Equal("foo bar baz\n", buf.String())
}

func TestTable_AddRow_Panic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic")
		}
	}()

	var buf bytes.Buffer
	New(&buf).AddRow("foo", "bar").AddRow("baz")
}

type Suite struct {
	suite.Suite
}

func (s *Suite) testTableRender(factory func(w io.Writer) *Table, expected string) {
	var buf bytes.Buffer

	// remove leading newline char
	expected = strings.TrimLeft(expected, "\n")
	expected = strings.ReplaceAll(expected, ".", " ")

	tab := factory(&buf)

	_, err := tab.Render()
	s.NoError(err)
	s.Equal(expected, buf.String())
}

func (s *Suite) TestTable_Render() {
	defer style.Enable()()

	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w)
		},
		``,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w).
				AddRow("foo\nbarbaz", "qux")
		},
		`
foo....qux
barbaz....
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithPadding(3)).
				AddRow("foo\nbarbaz", "qux")
		},
		`
foo......qux
barbaz......
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithPadding(-1)).
				AddRow("foo", "bar")
		},
		`
foobar
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(17), WithMargin(2)).
				AddRow("foo", "barbaz", "qux")
		},
		`
..foo.barb….qux..
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(10), WithMargin(-1)).
				AddRow("foo", "bar", "baz")
		},
		`
f….bar.baz
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(30)).
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(1, 2, "foo")
		},
		`
foo....bar....baz...
foofoo.barbar.bazbaz
1......2......foo...
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(40)).
				AddRow(
					"foo",
					bar.Bar{
						MaxWidth:       -1,
						Completed:      10,
						RemainingStyle: bar.NewStyle('-', nil),
						CompletedStyle: bar.NewStyle('#', nil),
						FinishedStyle:  bar.NewStyle('#', nil),
					},
					"qux",
				)
		},
		`
foo.###-----------------------------.qux
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(18)).
				AddRow(
					"foo",
					bar.Bar{
						MaxWidth:       -1,
						Completed:      10,
						RemainingStyle: bar.NewStyle('-', nil),
						CompletedStyle: bar.NewStyle('#', nil),
						FinishedStyle:  bar.NewStyle('#', nil),
					},
					bar.Bar{
						MaxWidth:       -1,
						Completed:      40,
						RemainingStyle: bar.NewStyle('-', nil),
						CompletedStyle: bar.NewStyle('#', nil),
						FinishedStyle:  bar.NewStyle('#', nil),
					},
					"qux",
				)
		},
		`
foo.----.##---.qux
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderAll)).
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(1, 2, "foo")
		},
		`
┌────────┬────────┬────────┐
│ foo    │ bar    │ baz    │
├────────┼────────┼────────┤
│ foofoo │ barbar │ bazbaz │
├────────┼────────┼────────┤
│ 1      │ 2      │ foo    │
└────────┴────────┴────────┘
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderRow)).
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(1, 2, "foo")
		},
		`
foo    bar    baz   
────────────────────
foofoo barbar bazbaz
────────────────────
1      2      foo   
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderRow|BorderColumn)).
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(1, 2, "foo")
		},
		`
foo    │ bar    │ baz   
───────┼────────┼───────
foofoo │ barbar │ bazbaz
───────┼────────┼───────
1      │ 2      │ foo   
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderAllVertical)).
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(1, 2, "foo")
		},
		`
│ foo    │ bar    │ baz    │
│ foofoo │ barbar │ bazbaz │
│ 1      │ 2      │ foo    │
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderAll^BorderRow)).
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(1, 2, "foo")
		},
		`
┌────────┬────────┬────────┐
│ foo    │ bar    │ baz    │
│ foofoo │ barbar │ bazbaz │
│ 1      │ 2      │ foo    │
└────────┴────────┴────────┘
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderAll), WithBorderStyle(style.New(style.FgBlack))).
				AddRow("foo", "bar")
		},
		"\x1b[30m┌\x1b[0m\x1b[30m─\x1b[0m\x1b[30m───\x1b[0m\x1b[30m─\x1b[0m\x1b[30m┬\x1b[0m\x1b[30m─\x1b[0m\x1b[30m───\x1b[0m\x1b[30m─\x1b[0m\x1b[30m┐\x1b[0m\n"+
			"\x1b[30m│\x1b[0m foo \x1b[30m│\x1b[0m bar \x1b[30m│\x1b[0m\n"+
			"\x1b[30m└\x1b[0m\x1b[30m─\x1b[0m\x1b[30m───\x1b[0m\x1b[30m─\x1b[0m\x1b[30m┴\x1b[0m\x1b[30m─\x1b[0m\x1b[30m───\x1b[0m\x1b[30m─\x1b[0m\x1b[30m┘\x1b[0m\n",
	)
}

func TestMain(t *testing.T) {
	suite.Run(t, new(Suite))
}
