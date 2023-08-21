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
