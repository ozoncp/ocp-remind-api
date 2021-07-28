package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemind_String(t *testing.T) {
	var tests = []struct {
		name     string
		input    Remind
		expected string
	}{
		{
			name: "Remind id=1, userid=23, date 01.01.1970, text='test_text'",
			input: Remind{
				Id:       1,
				UserId:   23,
				Deadline: time.Unix(0, 0),
				Text:     "test_text",
			},
			expected: "Remind\tId=1\tUserId=23\tDeadline=1970-01-01 03:00:00\tText=test_text",
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			assert.Equal(t, tCase.expected, tCase.input.String(), tCase.name)
		})
	}
}

func mustDuration(t *testing.T, str string) time.Duration {
	t.Helper()
	d, err := time.ParseDuration(str)
	require.NoError(t, err)
	return d
}

func mustTime(t *testing.T, str string) time.Time {
	t.Helper()
	tt, err := time.Parse("2006-01-02 15:04:05", str)
	require.NoError(t, err)
	return tt
}

func TestRemind_BeforeDeadline(t *testing.T) {
	var tests = []struct {
		name string
		now  time.Time
	}{
		{
			name: "deadline less than 1s in past",
			now:  time.Now(),
		},
	}
	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			r := NewRemind(1, 2, tCase.now, "test")
			duration := r.BeforeDeadline()
			assert.Less(t, duration, mustDuration(t, "0s"))
			assert.Greater(t, duration, mustDuration(t, "-1s"))
		})
	}
}

func TestRemind_Miss(t *testing.T) {
	var tests = []struct {
		name     string
		input    Remind
		expected bool
	}{
		{
			name:     "remind in the future, not missed ",
			input:    NewRemind(1, 1, mustTime(t, "2100-01-01 00:00:00"), "test"),
			expected: false,
		},

		{
			name:     "remind in the past, already missed",
			input:    NewRemind(1, 1, time.Now(), "test"),
			expected: true,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			assert.Equal(t, tCase.expected, tCase.input.Miss())
		})
	}
}

func TestRemind_Move(t *testing.T) {
	tests := []struct {
		name     string
		input    Remind
		duration time.Duration
		expected time.Time
	}{
		{
			name:     "move 2 days forward",
			input:    NewRemind(1, 2, mustTime(t, "2030-02-03 00:00:00"), "test"),
			duration: mustDuration(t, "48h"),
			expected: mustTime(t, "2030-02-05 00:00:00"),
		},
		{
			name:     "move 2 days backward",
			input:    NewRemind(1, 2, mustTime(t, "2030-02-03 00:00:00"), "test"),
			duration: mustDuration(t, "-48h"),
			expected: mustTime(t, "2030-02-01 00:00:00"),
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			tCase.input.Move(tCase.duration)
			assert.Equal(t, tCase.expected, tCase.input.Deadline)
		})
	}
}
