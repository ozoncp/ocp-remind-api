package repo

import (
	"github.com/ozoncp/ocp-remind-api/internal/models"
)

type RemindsRepo interface {
	Add(remind []models.Remind) error
	Describe(id uint64) (*models.Remind, error)
	List(limit, offset uint64) ([]models.Remind, error)
}
