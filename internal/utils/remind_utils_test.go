package utils

import (
	"errors"
	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	var tests = []struct {
		name  string
		input []models.Remind
		err   error
	}{
		{name: "three different reminds",
			input: []models.Remind{
				models.NewRemind(1, 1, time.Now(), "first"),
				models.NewRemind(2, 1, time.Now(), "second"),
				models.NewRemind(3, 1, time.Now(), "third")},
			err: nil,
		},
		{name: "reminds with same ids",
			input: []models.Remind{
				models.NewRemind(1, 1, time.Now(), "first"),
				models.NewRemind(1, 1, time.Now(), "second"),
				models.NewRemind(3, 1, time.Now(), "third")},
			err: errors.New("there are structs with same ids."),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := Map(test.input)
			assert.Equal(t, err, test.err)
			for k, v := range m {
				assert.Equal(t, k, v.Id)
			}

		})
	}
}


func TestBatchReminds(t *testing.T) {
	var tests = []struct {
		name string
		input     []models.Remind
		batchSize int
		expected  [][]models.Remind
	}{
		{
			name: "simple batch",
			input: []models.Remind{
			models.NewRemind(1, 4, time.Now(), "first"),
			models.NewRemind(2, 5, time.Now(), "second"),
			models.NewRemind(3, 6, time.Now(), "third"),
		},
		batchSize: 2,
		expected: [][]models.Remind{
			{
				models.NewRemind(1, 4, time.Now(), "first"),
				models.NewRemind(2, 5, time.Now(), "second"),
			}, {
				models.NewRemind(3, 6, time.Now(), "third"),
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t* testing.T){
			for i, v:= range BatchReminds(test.input, test.batchSize){
				assert.Equal(t, len(v), len(test.expected[i]))
				for j, r := range v{
					assert.Equal(t, r.Text, test.expected[i][j].Text)
					assert.Equal(t, r.Id, test.expected[i][j].Id)
					assert.Equal(t, r.UserId, test.expected[i][j].UserId)
				}
			}
		})
	}
}
