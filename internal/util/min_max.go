package util

// MinInt returns the smaller of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// MaxInt returns the bigger of a and b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinFloat64 returns the smaller of a and b.
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

// MaxFloat64 returns the bigger of a and b.
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}

	return b
}
