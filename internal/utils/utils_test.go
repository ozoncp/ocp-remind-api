package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMirror(t *testing.T) {
	var tests = []struct {
		input    map[string]int
		expected map[int]string
	}{
		{map[string]int{"one": 1, "two": 2, "three": 3},
			map[int]string{1: "one", 2: "two", 3: "three"},
		},
		{map[string]int{"Russia": 3, "Finland": 3, "France": 5},
			map[int]string{3: "Finland", 5: "France"},
		},
	}

	for _, test := range tests {
		result := Mirror(test.input)
		assert.True(t, reflect.DeepEqual(result, test.expected))
	}

	require.NotPanics(t, func() { var nilMap map[string]int; Mirror(nilMap) },
		"func panic on input is nil")
}

func TestBatch(t *testing.T) {
	var tests = []struct {
		input     []int
		batchSize int
		expected  [][]int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			3,
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}},
		{[]int{1, 2, 3},
			1,
			[][]int{{1}, {2}, {3}}},
	}
	for _, test := range tests {
		var result = Batch(test.input, test.batchSize)
		assert.Equal(t, result, test.expected)
	}
}

func TestFilter(t *testing.T) {
	var tests = []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			[]int{1, 2, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20}},
		{[]int{},
			[]int{}},
		{[]int{100, 200, 300},
			[]int{100, 200, 300}},
	}

	for _, test := range tests {
		assert.Equal(t, Filter(test.input), test.expected)
	}
}
