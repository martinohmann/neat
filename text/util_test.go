package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const lorem = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."

func TestCountLines(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(1, CountLines(""))
	assert.Equal(1, CountLines("foo"))
	assert.Equal(2, CountLines("foo\nbar"))
	assert.Equal(3, CountLines("foo\nbar\n"))
	assert.Equal(3, CountLines("foo\nbar\nbaz"))
}

func TestDisplayWidth(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(0, DisplayWidth(""))
	assert.Equal(3, DisplayWidth("\nfoo"))
	assert.Equal(6, DisplayWidth("foo\nbarbaz"))
}

func TestWrapWords(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(`Lorem ipsum dolor sit amet,
consetetur sadipscing elitr,
sed diam nonumy eirmod tempor
invidunt ut labore et dolore
magna aliquyam erat, sed diam
voluptua. At vero eos et
accusam et justo duo dolores
et ea rebum. Stet clita kasd
gubergren, no sea takimata
sanctus est Lorem ipsum dolor
sit amet.`,
		WrapWords(lorem, 30),
	)

	assert.Equal(`Lorem
ipsum
dolor sit
amet,
consetetur
sadipscing
elitr,
sed diam
nonumy
eirmod
tempor
invidunt
ut labore
et dolore
magna
aliquyam
erat, sed
diam
voluptua.
At vero
eos et
accusam
et justo
duo
dolores
et ea
rebum.
Stet
clita
kasd
gubergren,
no sea
takimata
sanctus
est Lorem
ipsum
dolor sit
amet.`,
		WrapWords(lorem, 9),
	)
}
