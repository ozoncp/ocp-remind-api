package saver

import (
	"time"

	"github.com/ozoncp/ocp-remind-api/internal/flusher"
	"github.com/ozoncp/ocp-remind-api/internal/models"
)

type Saver interface {
	Save(remind models.Remind)
	Init()
	Close()
}

type remindSaver struct {
	flusher flusher.Flusher
	period  time.Duration
	bufC    chan models.Remind
	closeC  chan struct{}
	reminds []models.Remind
}

func (rs *remindSaver) Init() {
	go func() {
		ticker := time.NewTicker(rs.period)
		defer ticker.Stop()

		for {
			select {
			case remind := <-rs.bufC:
				rs.reminds = append(rs.reminds, remind)
			case <-ticker.C:
				rs.reminds = rs.flusher.Flush(rs.reminds)
			case <-rs.closeC:
				if len(rs.reminds) != 0 {
					rs.flusher.Flush(rs.reminds)
				}

				return
			}
		}
	}()
}

func (rs remindSaver) Save(remind models.Remind) {
	rs.bufC <- remind
}

func (rs remindSaver) Close() {
	rs.closeC <- struct{}{}
	close(rs.bufC)
	close(rs.closeC)
}

type Option func(*remindSaver)

func WithDuration(v time.Duration) Option {
	return func(s *remindSaver) {
		if v > 0 {
			s.period = v
		}
	}
}

//This function in addition to creating object,
//also initialize it (call Init() method)
func NewSaver(capacity uint, flusher flusher.Flusher, opts ...Option) Saver {
	saver := remindSaver{
		flusher: flusher,
		period:  5 * time.Second,
		reminds: make([]models.Remind, 0, capacity),
		bufC:    make(chan models.Remind),
		closeC:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(&saver)
	}

	saver.Init()

	return &saver
}
