package models

type timestamp int64

type Remind struct {
	Id       uint64
	UserId   uint64
	Deadline timestamp
	Text     string
}
