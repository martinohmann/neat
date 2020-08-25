package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStyleString(t *testing.T) {
	defer Enable()()
	assert := assert.New(t)

	assert.Equal("", StyleString(""))
	assert.Equal("string", StyleString("string"))
	assert.Equal("{}string", StyleString("{}string"))
	assert.Equal("{foo\x1b[31mstring\x1b[0m", StyleString("{foo{red}string{reset}"))
	assert.Equal("\x1b[31;1mred\x1b[0m", StyleString("{red,bold}red{reset}"))
	assert.Equal("{unknown}string\x1b[42;30mwithbg\x1b[0m", StyleString("{unknown}string{bggreen,black}withbg{reset}"))
	assert.Equal("{red,unknown}string", StyleString("{red,unknown}string"))
}
