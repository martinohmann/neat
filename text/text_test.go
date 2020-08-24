package text

import (
	"testing"

	"github.com/martinohmann/neat/measure"
	"github.com/stretchr/testify/assert"
)

func TestText_Render(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("   ", New("").Render(3))
	assert.Equal("foo", New("foo").Render(3))
	assert.Equal("fo"+string(ellipsis), New("foobar").Render(3))
	assert.Equal("foob"+string(ellipsis), New("foobar").Render(5))
	assert.Equal("foobar  ", New("foobar").Render(8))
}

func TestText_Render_WordWrap(t *testing.T) {
	assert := assert.New(t)

	text := Text{
		Text:      lorem,
		WordWrap:  true,
		Alignment: AlignJustify,
	}

	assert.Equal(`Lorem ipsum dolor sit    amet,
consetetur  sadipscing  elitr,
sed diam nonumy eirmod  tempor
invidunt ut labore et   dolore
magna aliquyam erat, sed  diam
voluptua.  At  vero  eos    et
accusam et justo duo   dolores
et ea rebum. Stet clita   kasd
gubergren,  no  sea   takimata
sanctus est Lorem ipsum  dolor
sit                      amet.`,
		text.Render(30),
	)
}

func TestBar_Measure(t *testing.T) {
	assert := assert.New(t)

	mm := measure.NewMeasurement

	assert.Equal(mm(0, 0), New("").Measure(10))
	assert.Equal(mm(3, 3), New("foo").Measure(10))
	assert.Equal(mm(3, 3), New("foobar").Measure(3))
}
