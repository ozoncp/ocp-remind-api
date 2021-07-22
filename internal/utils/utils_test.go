package utils

import (
	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMirror(t *testing.T) {
	source := make(map[string]int)
	source["one"] = 1
	source["two"] = 2
	source["three"] = 3
	mirrored := Mirror(source)
	assert.Equal(t, len(mirrored), len(source), "map size should be the same")
	assert.Contains(t, mirrored, 1, "map should contains this key")
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
			[]int{1, 2, 20, 4, 18, 6, 16, 8, 9, 10, 15, 12, 14}},
		{[]int{},
			[]int{}},
		{[]int{100, 200, 300},
			[]int{100, 200, 300}},
	}

	for _, test := range tests {
		assert.Equal(t, Filter(test.input), test.expected)
	}
}

func TestToMap(t *testing.T) {
	var tests = []struct {
		input    []models.Remind
		expected map[uint64]models.Remind
	}{
		{[]models.Remind{{1, 1, 1000, "first"},
			{2, 2, 2000, "second"},
			{3, 3, 3000, "third"}},
			map[uint64]models.Remind{1: {1, 1, 1000, "first"},
				2: {2, 2, 2000, "second"},
				3: {3, 3, 3000, "third"}}},
	}
	for _, test := range tests {
		result, _ := ToMap(test.input)
		assert.Equal(t, result, test.expected)

	}
}
