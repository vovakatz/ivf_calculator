package utils

import "strconv"

// ParseFloat is a helper function to parse float values
// while ignoring error
func ParseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
