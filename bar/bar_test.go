package bar

import (
	"testing"

	"github.com/martinohmann/neat/measure"
	"github.com/stretchr/testify/assert"
)

func TestBar_Render(t *testing.T) {
	assert := assert.New(t)

	bar := Bar{
		RemainingStyle: NewStyle('r', nil),
		CompletedStyle: NewStyle('c', nil),
		FinishedStyle:  NewStyle('f', nil),
		Completed:      50,
	}

	assert.Equal("ccrr", bar.Render(4))
	assert.Equal("ccrrr", bar.Render(5))
	bar.Completed = 66.6
	assert.Equal("ccccccrrrr", bar.Render(10))
	bar.Completed = 100
	assert.Equal("ffffffffff", bar.Render(10))
}

func TestBar_Measure(t *testing.T) {
	assert := assert.New(t)

	measureBar := func(bw, mw int) measure.Measurement {
		b := Bar{MaxWidth: bw}
		return b.Measure(mw)
	}

	mm := measure.NewMeasurement

	assert.Equal(mm(4, 4), measureBar(0, 10))
	assert.Equal(mm(4, 5), measureBar(5, 10))
	assert.Equal(mm(4, 10), measureBar(-1, 10))
	assert.Equal(mm(4, 10), measureBar(20, 10))
}
