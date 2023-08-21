package util

import "strconv"

// ToInt convert string to int
func ToInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}
