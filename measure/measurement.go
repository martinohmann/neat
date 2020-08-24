package measure

import (
	"github.com/martinohmann/neat/internal/util"
)

// Measurement is a pair of minimum and maximum values that can be used to
// express requested widths and heights for console.Renderable objects.
type Measurement struct {
	Minimum int
	Maximum int
}

// NewMeasurement creates a new Measurement with min and max values.
func NewMeasurement(min, max int) Measurement {
	return Measurement{
		Minimum: min,
		Maximum: max,
	}
}

// Normalize ensures that minimum is always >= 0 and that maximum is always >=
// minimum. Returns a new Measurement.
func (m Measurement) Normalize() Measurement {
	minimum := util.MinInt(util.MaxInt(0, m.Minimum), m.Maximum)

	return NewMeasurement(
		util.MaxInt(0, minimum),
		util.MaxInt(0, util.MaxInt(minimum, m.Maximum)),
	)
}

// Span returns the span of the Measurement. E.g. if min is 3 and max is 10
// then span will be 7.
func (m Measurement) Span() int {
	m = m.Normalize()
	return m.Maximum - m.Minimum
}
