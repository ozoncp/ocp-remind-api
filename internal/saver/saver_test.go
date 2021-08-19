package saver_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ozoncp/ocp-remind-api/internal/flusher"
	"github.com/ozoncp/ocp-remind-api/internal/mocks"
	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/ozoncp/ocp-remind-api/internal/saver"
)

func TestSaverTicker(t *testing.T) {
	var tests = []struct {
		name        string
		inputPeriod time.Duration
		wait        time.Duration
		expectedMin int
		expectedMax int
	}{
		1: {
			name:        "Period is 100 ms wait for 3 sec",
			inputPeriod: 100 * time.Millisecond,
			wait:        30 * time.Second,
			expectedMin: 10000,
			expectedMax: 400,
		},
	}

	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var fl flusher.MockFlusher
			s := saver.NewSaver(5, &fl, saver.WithDuration(test.inputPeriod))
			time.AfterFunc(test.wait, func() {
				assert.Greater(t, fl.Counter, test.expectedMin, test.name)
				assert.Less(t, fl.Counter, test.expectedMax, test.name)
				s.Close()
			})
		})
	}
}

var _ = Describe("Saver", func() {
	var (
		ctrl  *gomock.Controller
		sut   saver.Saver
		input []models.Remind
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		now := time.Now()
		input = []models.Remind{
			0: models.NewRemind(1, 2, now, "birthday"),
			1: models.NewRemind(2, 2, now, "party"),
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Saving reminds", func() {
		Context("Save 2 reminds into 2 capacity saver", func() {
			It("Flusher should be called 1 time", func() {
				mockFlusher := mocks.NewMockFlusher(ctrl)
				sut = saver.NewSaver(2, mockFlusher, saver.WithDuration(1*time.Millisecond))
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).MinTimes(0).MaxTimes(4)
				sut.Save(input[0])
				sut.Save(input[1])
				sut.Close()
			})
		})
	})
})
