package utils

import (
	"fmt"

	"go.uber.org/multierr"

	"github.com/ozoncp/ocp-remind-api/internal/models"
)

func Map(reminds []models.Remind) (map[uint64]models.Remind, []error) {
	result := make(map[uint64]models.Remind, len(reminds))
	var err error
	for _, remind := range reminds {
		if _, contains := result[remind.Id]; !contains {
			result[remind.Id] = remind
		} else {
			err = multierr.Append(err, fmt.Errorf("%d: already persists in as a reminder", remind.Id))
		}
	}
	return result, multierr.Errors(err)
}

func BatchReminds(input []models.Remind, size int) [][]models.Remind {
	if len(input) <= size {
		return [][]models.Remind{input}
	}

	result := make([][]models.Remind, 0, (len(input)+size-1)/size)
	for i, j := 0, size; j <= len(input); i, j = i+size, j+size {
		result = append(result, input[i:j])
	}
	if v := len(input) % size; v != 0 {
		result = append(result, input[len(input)-v:])
	}

	return result
}
