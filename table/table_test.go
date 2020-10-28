package table

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/martinohmann/neat/bar"
	"github.com/martinohmann/neat/style"
	"github.com/martinohmann/neat/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const lorem = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."

func TestTable_Render_Reset(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	table := New(&buf)
	table.AddRow("foo\nbarbaz", "qux").AddRow("one", "two")

	nlines, err := table.RenderLines()
	assert.NoError(err)
	assert.Equal(3, nlines)
	assert.Equal("foo    qux\nbarbaz    \none    two\n", buf.String())

	buf.Reset()
	table.Reset()

	table.AddRow("foo", "bar", "baz")

	nlines, err = table.RenderLines()
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

func TestTable_AddRow_PanicOrder(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected panic")
		}
	}()

	var buf bytes.Buffer
	New(&buf).AddRow("foo", "bar").AddHeader("baz", "qux")
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

	s.NoError(tab.Render())
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
			return New(w, WithBorderMask(BorderAll)).
				AddHeader("foo", "bar", "baz").
				AddHeader("one", "two", "three").
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(12, 13, 14).
				AddFooter(1, 2, "foo")
		},
		`
┌────────┬────────┬────────┐
│ foo    │ bar    │ baz    │
╞════════╪════════╪════════╡
│ one    │ two    │ three  │
╞════════╪════════╪════════╡
│ foo    │ bar    │ baz    │
├────────┼────────┼────────┤
│ foofoo │ barbar │ bazbaz │
├────────┼────────┼────────┤
│ 12     │ 13     │ 14     │
╞════════╪════════╪════════╡
│ 1      │ 2      │ foo    │
└────────┴────────┴────────┘
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithBorderMask(BorderSection|BorderColumn)).
				AddHeader("foo", "bar", "baz").
				AddHeader("one", "two", "three").
				AddRow("foo", "bar", "baz").
				AddRow("foofoo", "barbar", "bazbaz").
				AddRow(12, 13, 14).
				AddFooter(1, 2, "foo")
		},
		`
foo    │ bar    │ baz   
═══════╪════════╪═══════
one    │ two    │ three 
═══════╪════════╪═══════
foo    │ bar    │ baz   
foofoo │ barbar │ bazbaz
12     │ 13     │ 14    
═══════╪════════╪═══════
1      │ 2      │ foo   
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
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(14), WithWordWrap(true)).
				AddRow("foo bar", "baz qux")
		},
		`
foo....baz.qux
bar...........
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w, WithMaxWidth(12), WithColumnWordWrap(false, true)).
				AddRow("foo bar", "baz qux")
		},
		`
foo.….baz...
......qux...
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w,
				WithBorderMask(BorderAll),
				WithWordWrap(true),
				WithAlignment(text.AlignRight)).
				AddRow(lorem[:75], lorem[:75])
		},
		`
┌──────────────────────────────────────┬───────────────────────────────────────┐
│          Lorem ipsum dolor sit amet, │           Lorem ipsum dolor sit amet, │
│     consetetur sadipscing elitr, sed │ consetetur sadipscing elitr, sed diam │
│                       diam nonumy ei │                             nonumy ei │
└──────────────────────────────────────┴───────────────────────────────────────┘
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w,
				WithBorderMask(BorderAll),
				WithWordWrap(true),
				WithColumnAlignment(text.AlignLeft, text.AlignJustify, text.AlignRight, text.AlignCenter)).
				AddRow(lorem[:75], lorem[:75], lorem[:75], lorem[:75])
		},
		`
┌──────────────────┬───────────────────┬───────────────────┬───────────────────┐
│ Lorem ipsum      │ Lorem ipsum dolor │ Lorem ipsum dolor │ Lorem ipsum dolor │
│ dolor sit amet,  │ sit         amet, │         sit amet, │     sit amet,     │
│ consetetur       │ consetetur        │        consetetur │    consetetur     │
│ sadipscing       │ sadipscing elitr, │ sadipscing elitr, │ sadipscing elitr, │
│ elitr, sed diam  │ sed  diam  nonumy │   sed diam nonumy │  sed diam nonumy  │
│ nonumy ei        │ ei                │                ei │        ei         │
└──────────────────┴───────────────────┴───────────────────┴───────────────────┘
`,
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w,
				WithStyle(style.New(style.Bold))).
				AddRow("foo", "bar")
		},
		"\x1b[1mfoo\x1b[0m \x1b[1mbar\x1b[0m\n",
	)
	s.testTableRender(
		func(w io.Writer) *Table {
			return New(w,
				WithColumnStyle(style.New(style.Bold))).
				AddRow("foo", "bar")
		},
		"\x1b[1mfoo\x1b[0m bar\n",
	)
}

func TestMain(t *testing.T) {
	suite.Run(t, new(Suite))
}
