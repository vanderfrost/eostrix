package utils

// min for compat with go version before 1.25
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
