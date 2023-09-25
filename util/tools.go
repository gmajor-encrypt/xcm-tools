package util

import "strconv"

// ToInt convert string to int
func ToInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

func ToUint(i string) uint {
	if i, err := strconv.ParseUint(i, 10, 0); err == nil {
		return uint(i)
	}
	return 0
}
