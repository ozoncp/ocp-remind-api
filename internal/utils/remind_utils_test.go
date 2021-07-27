package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	now := time.Now()
	var tests = []struct {
		name  string
		input []models.Remind
		err   []error
	}{
		{
			name: "3 different reminds",
			input: []models.Remind{
				models.NewRemind(1, 1, now, "first"),
				models.NewRemind(2, 1, now, "second"),
				models.NewRemind(3, 1, now, "third"),
			},
			err: nil,
		},
		{
			name: "total 3 reminds, 2 reminds with same ids",
			input: []models.Remind{
				models.NewRemind(1, 1, now, "first"),
				models.NewRemind(1, 1, now, "second"),
				models.NewRemind(3, 1, now, "third"),
			},
			err: []error{errors.New("1: already persists in as a reminder")},
		},
		{
			name: "reminds with same ids, 2 pairs in map and 2 errors expected",
			input: []models.Remind{
				models.NewRemind(1, 1, now, "first"),
				models.NewRemind(1, 1, now, "second"),
				models.NewRemind(3, 1, now, "third"),
				models.NewRemind(3, 1, now, "third"),
			},
			err: []error{errors.New("1: already persists in as a reminder"),
				errors.New("3: already persists in as a reminder"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := Map(test.input)
			assert.Equal(t, test.err, err)
			for k, v := range m {
				assert.Equal(t, v.Id, k)
			}

		})
	}
}

func TestBatchReminds(t *testing.T) {
	now := time.Now()
	var tests = []struct {
		name      string
		input     []models.Remind
		batchSize int
		expected  [][]models.Remind
	}{
		{
			name: "batch 3 reminds by 2, 2 + 1 reminds expected",
			input: []models.Remind{
				models.NewRemind(1, 4, now, "first"),
				models.NewRemind(2, 5, now, "second"),
				models.NewRemind(3, 6, now, "third"),
			},
			batchSize: 2,
			expected: [][]models.Remind{
				{
					models.NewRemind(1, 4, now, "first"),
					models.NewRemind(2, 5, now, "second"),
				},
				{
					models.NewRemind(3, 6, now, "third"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i, v := range BatchReminds(test.input, test.batchSize) {
				assert.Equal(t, len(v), len(test.expected[i]))
				for j, r := range v {
					assert.Equal(t, test.expected[i][j].Text, r.Text)
					assert.Equal(t, test.expected[i][j].Id, r.Id)
					assert.Equal(t, test.expected[i][j].UserId, r.UserId)
				}
			}
		})
	}
}
