package utils

// OnlyOneTrue returns true if only 1 value out of 3 is set to true.
func OnlyOneTrue(b1, b2, b3 bool) bool {
	count := 0
	if b1 {
		count++
	}
	if b2 {
		count++
	}
	if b3 {
		count++
	}
	return count == 1
}
