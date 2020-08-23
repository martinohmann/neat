package util

func CountDigitsInt64(num int64) (n int) {
	for num != 0 {
		num /= 10
		n += 1
	}

	return n
}
