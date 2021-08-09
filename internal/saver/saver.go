package saver

import (
	"errors"
	"time"

	"github.com/ozoncp/ocp-remind-api/internal/flusher"
	"github.com/ozoncp/ocp-remind-api/internal/models"
)

type Saver interface {
	Save(remind models.Remind) error
	Init()
	Close()
}

type remindSaver struct {
	flusher   flusher.Flusher
	ticker    *time.Ticker
	reminds   []models.Remind
	saveChan  chan models.Remind
	closeChan chan struct{}
}

func (saver *remindSaver) Init() {
	go func() {
		defer saver.ticker.Stop()

		for {
			select {
			case remind := <-saver.saveChan:
				saver.reminds = append(saver.reminds, remind)
			case <-saver.ticker.C:
				if len(saver.reminds) != 0 {
					saver.reminds = saver.flusher.Flush(saver.reminds)
				}
			case <-saver.closeChan:
				if len(saver.reminds) != 0 {
					saver.flusher.Flush(saver.reminds)
				}

				close(saver.saveChan)
				close(saver.closeChan)

				return
			}
		}
	}()
}

var ErrCapacity = errors.New("there is no avaible capacity")

func (saver remindSaver) Save(remind models.Remind) error {
	if len(saver.reminds) < cap(saver.reminds) {
		saver.saveChan <- remind

		return nil
	}

	return ErrCapacity
}

func (saver remindSaver) Close() {
	saver.closeChan <- struct{}{}
}

func NewSaver(capacity uint, flusher flusher.Flusher, duration ...time.Duration) Saver {
	d := 5 * time.Second
	if len(duration) > 0 {
		d = duration[0]
	}

	saver := remindSaver{
		flusher:   flusher,
		ticker:    time.NewTicker(d),
		reminds:   make([]models.Remind, 0, capacity),
		saveChan:  make(chan models.Remind),
		closeChan: make(chan struct{}),
	}
	saver.Init()

	return &saver
}
