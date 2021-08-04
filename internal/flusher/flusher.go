package flusher

import (
	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/ozoncp/ocp-remind-api/internal/repo"
	"github.com/ozoncp/ocp-remind-api/internal/utils"
)

type Flusher interface {
	Flush([]models.Remind) []models.Remind
}

type remindsFlusher struct {
	repo      repo.RemindsRepo
	chunkSize int
}

func (f *remindsFlusher) Flush(reminds []models.Remind) []models.Remind {
	batched := utils.BatchReminds(reminds, f.chunkSize)
	notAdded := make([]models.Remind, 0, len(reminds))
	for _, v := range batched {
		err := f.repo.Add(v)
		if err != nil {
			notAdded = append(notAdded, v...)
		}
	}
	return notAdded
}

func NewFlusher(r repo.RemindsRepo, chunkSize int) Flusher {
	return &remindsFlusher{repo: r, chunkSize: chunkSize}
}
