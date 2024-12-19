package utils

// Contains returns true if value can be found in the array.
// Otherwise, false.
func Contains(arr []string, value string) bool {
	for _, item := range arr {
		if item == value {
			return true
		}
	}
	return false
}
