package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMirror(t *testing.T) {
	var tests = []struct {
		name     string
		input    map[string]int
		expected map[int]string
	}{
		{
			name:     "three different pairs",
			input:    map[string]int{"one": 1, "two": 2, "three": 3},
			expected: map[int]string{1: "one", 2: "two", 3: "three"},
		},
		{
			name:     "three pairs, two are same",
			input:    map[string]int{"Russia": 3, "Finland": 3, "France": 5},
			expected: map[int]string{3: "Finland", 5: "France"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.True(t, reflect.DeepEqual(test.expected, Mirror(test.input)))
		})
	}

	require.NotPanics(t, func() { var nilMap map[string]int; Mirror(nilMap) },
		"function panics on nil input")
}

func TestBatch(t *testing.T) {
	var tests = []struct {
		name      string
		input     []int
		batchSize int
		expected  [][]int
	}{
		{
			name:      "ten elements batched by three elements",
			input:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			batchSize: 3,
			expected:  [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}},
		{
			name:      "three elements batched by one",
			input:     []int{1, 2, 3},
			batchSize: 1,
			expected:  [][]int{{1}, {2}, {3}},
		},
	}
	for _, test := range tests {
		var result = Batch(test.input, test.batchSize)
		assert.Equal(t, test.expected, result)
	}
}

func TestFilter(t *testing.T) {
	var tests = []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "from 1 to 20 filter prime numbers",
			input:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			expected: []int{1, 2, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20},
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "full slice goes throw the filter",
			input:    []int{100, 200, 300},
			expected: []int{100, 200, 300},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, Filter(test.input))
	}
}
