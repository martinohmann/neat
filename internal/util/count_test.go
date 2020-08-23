package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountDigitsInt64(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(19, CountDigitsInt64(math.MinInt64))
	assert.Equal(5, CountDigitsInt64(-99999))
	assert.Equal(2, CountDigitsInt64(-10))
	assert.Equal(1, CountDigitsInt64(-9))
	assert.Equal(1, CountDigitsInt64(-1))
	assert.Equal(1, CountDigitsInt64(9))
	assert.Equal(2, CountDigitsInt64(10))
	assert.Equal(5, CountDigitsInt64(99999))
	assert.Equal(19, CountDigitsInt64(math.MaxInt64))
}
