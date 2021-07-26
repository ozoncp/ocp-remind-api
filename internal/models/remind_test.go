package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestRemind_String(t *testing.T) {
	var tests = []struct{
		name string
		input Remind
		expected string
	}{
		{
			name:  "Remind id=1, userid=23, date 01.01.1970, text='test_text'",
			input: Remind{1,23,timestamp(time.Unix(0,0)), "test_text"},
			expected: "Remind\tId=1\tUserId=23\tDeadline=1970-01-01 03:00:00\tText=test_text",
		},
	}

	for _, tcase := range tests{
		t.Run(tcase.name, func(t *testing.T) {
			assert.Equal(t, tcase.input.String(), tcase.expected, tcase.name)
		})
	}
}

func TestNewRemind(t *testing.T) {
	var tests = []struct{
		name string
		id uint64
		userId uint64
		t time.Time
		text string
		expected Remind
	}{
		{
			name:   "Simple test",
			id:     1,
			userId: 23,
			t:      time.Unix(4,5),
			text:   "test_text",
			expected: Remind{1, 23, timestamp(time.Unix(4,5)), "test_text"},
		},
	}

	for _, tcase := range tests{
		t.Run(tcase.name, func(t* testing.T){
			r := NewRemind(tcase.id, tcase.userId, tcase.t, tcase.text)
			assert.Equal(t, r, tcase.expected)
		})
	}
}

func parseDurationWithoutError(str string)time.Duration{
	duration, err := time.ParseDuration(str)
	if err != nil{
		panic("error on parse duration without error")
	}
	return duration
}

func ParseTimeWithoutError(str string)time.Time{
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil{
		panic("error on parse time without error")
	}
	return t
}

func TestRemind_BeforeDeadline(t *testing.T) {
	var tests = []struct{
		name string
		now time.Time
	}{
		{
			name: "remind always in past",
			now: time.Now(),
		},
	}
	for _, tcase := range tests{
		t.Run(tcase.name, func(t* testing.T){
			r:= NewRemind(1,2,tcase.now, "test")
			duration := r.BeforeDeadline()
			assert.Less(t, duration, parseDurationWithoutError("0s"))
			assert.Greater(t, duration, parseDurationWithoutError("-1s"))
		})
	}
}

func TestRemind_Miss(t *testing.T) {
	var tests = []struct{
		name string
		input Remind
		expected bool
	}{
		{
			name: "remind in the future",
			input: NewRemind(1,1, ParseTimeWithoutError("2100-01-01 00:00:00"), "test"),
			expected: false,
		},

		{
			name: "remind in the past",
			input: NewRemind(1,1, time.Now(), "test"),
			expected: true,
		},
	}

	for _, tcase:= range tests{
		t.Run(tcase.name, func(t* testing.T){
			assert.Equal(t, tcase.input.Miss(), tcase.expected)
		})
	}
}

func TestRemind_Move(t *testing.T) {
	tests:= []struct{
		name string
		input Remind
		duration time.Duration
		expected timestamp
	}{
		{
			name: "move 2 days forward",
			input: NewRemind(1,2, ParseTimeWithoutError("2030-02-03 00:00:00"), "test"),
			duration: parseDurationWithoutError("48h"),
			expected: timestamp(ParseTimeWithoutError("2030-02-05 00:00:00")),
		},
		{
			name: "move 2 days backward",
			input: NewRemind(1,2, ParseTimeWithoutError("2030-02-03 00:00:00"), "test"),
			duration: parseDurationWithoutError("-48h"),
			expected: timestamp(ParseTimeWithoutError("2030-02-01 00:00:00")),
		},
	}

	for _, tcase:= range tests{
		t.Run(tcase.name, func(t* testing.T){
			tcase.input.Move(tcase.duration)
			assert.Equal(t, tcase.input.Deadline, tcase.expected)
		})
	}
}