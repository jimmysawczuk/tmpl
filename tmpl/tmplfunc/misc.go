package tmplfunc

// Seq returns a list of integers between 0 and max. If max is 0,
// nil is returned.
func Seq(max int) []int {
	if max < 0 {
		return nil
	}

	v := make([]int, 0, max)
	for i := 0; i < max; i++ {
		v = append(v, i)
	}

	return v
}

// Add returns the sum of the two arguments.
func Add(a, b int) int {
	return a + b
}

// Sub returns the difference of the two arguments.
func Sub(a, b int) int {
	return a - b
}
