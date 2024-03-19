package util

import "testing"

func Test_ToInt(t *testing.T) {
	cases := []struct {
		s        string
		expected int
	}{
		{s: "999999999", expected: 999999999},
		{s: "22222222", expected: 22222222},
		{s: "fff", expected: 0},
		{s: "-1123123123", expected: -1123123123},
	}
	for _, v := range cases {
		if ToInt(v.s) != v.expected {
			t.Errorf("ToInt(%s) != %d", v.s, v.expected)
		}
	}
}

func Test_HexToUint64(t *testing.T) {
	cases := []struct {
		s        string
		expected uint64
	}{
		{s: "0x123", expected: 291},
		{s: "0x1234567890", expected: 78187493520},
		{s: "0x1234567890abcdef", expected: 1311768467294899695},
		{s: "0x1234567890abcdef1234567890abcdef", expected: 0},
	}
	for _, v := range cases {
		if got := HexToUint64(v.s); got != v.expected {
			t.Errorf("HexToUint64(%s) != %d, got %d", v.s, v.expected, got)
		}
	}
}

func Test_ToUint(t *testing.T) {
	cases := []struct {
		s        string
		expected uint
	}{
		{s: "999999999", expected: 999999999},
		{s: "22222222", expected: 22222222},
		{s: "fff", expected: 0},
		{s: "-1123123123", expected: 0},
	}
	for _, v := range cases {
		if ToUint(v.s) != v.expected {
			t.Errorf("ToUint(%s) != %d", v.s, v.expected)
		}
	}
}
