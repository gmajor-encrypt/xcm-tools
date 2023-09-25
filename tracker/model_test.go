package tracker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Enum(t *testing.T) {
	cases := []struct {
		in   Enum
		want interface{}
	}{
		{Enum{"Error": "test"}, "test"},
		{Enum{"Error": 1}, 1},
		{Enum{"Error": []string{"test"}}, []string{"test"}},
		{Enum{"Error": []int{1}}, []int{1}},
		{Enum{"Error": []interface{}{1}}, []interface{}{1}},
		{Enum{"Error": []interface{}{1, "test"}}, []interface{}{1, "test"}},
	}
	for _, c := range cases {
		assert.Equal(t, c.in.Value(), c.want)
	}
}
