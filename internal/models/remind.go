package models

import (
	"fmt"
	"time"
)

type Remind struct {
	Deadline time.Time
	Text     string
	Id       uint64
	UserId   uint64
}

func NewRemind(id, userId uint64, time time.Time, text string) Remind {
	return Remind{
		Deadline: time,
		Text:     text,
		Id:       id,
		UserId:   userId,
	}
}

func (r Remind) BeforeDeadline() time.Duration {
	return time.Until(r.Deadline)
}

func (r Remind) Miss() bool {
	return time.Now().After(r.Deadline)
}

func (r *Remind) Move(duration time.Duration) {
	r.Deadline = r.Deadline.Add(duration)
}

func (r Remind) String() string {
	return fmt.Sprintf("Remind\tId=%d\tUserId=%d\tDeadline=%s\tText=%s",
		r.Id,
		r.UserId,
		r.Deadline.Format("2006-01-02 15:04:05"),
		r.Text,
	)
}
