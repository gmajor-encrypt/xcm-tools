package util

import (
	"github.com/itering/scale.go/utiles"
	"strconv"
)

// HexToUint64 convert hex string to uint64
func HexToUint64(h string) uint64 {
	blockNum, err := strconv.ParseUint(utiles.TrimHex(h), 16, 64)
	if err != nil {
		return 0
	}
	return blockNum
}

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

func InSlice[T comparable](el T, list []T) bool {
	for _, v := range list {
		if el == v {
			return true
		}
	}
	return false
}
