package models

import (
	"strconv"
	"time"
)

type timestamp time.Time


type Remind struct {
	Id       uint64
	UserId   uint64
	Deadline timestamp
	Text     string
}

func NewRemind(id, userId uint64, time time.Time, text string) Remind {
	return Remind{id,
		userId,
		timestamp(time),
		text}
}

func (r *Remind) BeforeDeadline() time.Duration {
	return time.Time(r.Deadline).Sub(time.Now())
}

func (r *Remind) Miss() bool {
	return time.Now().After(time.Time(r.Deadline))
}

func (r *Remind) Move(duration time.Duration) {
	r.Deadline = timestamp(time.Time(r.Deadline).Add(duration))
}

func (r *Remind) String() string {
	return "Remind\tId=" + strconv.FormatUint(r.Id, 10) +
		"\tUserId=" + strconv.FormatUint(r.UserId, 10) +
		"\tDeadline=" + time.Time(r.Deadline).Format("2006-01-02 15:04:05") +
		"\tText=" + r.Text
}
