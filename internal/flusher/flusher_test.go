package flusher_test

import (
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/ozoncp/ocp-remind-api/internal/flusher"
	"github.com/ozoncp/ocp-remind-api/internal/mocks"
	"github.com/ozoncp/ocp-remind-api/internal/models"
)

var _ = Describe("Flusher", func() {
	var (
		mr        *mocks.MockRepo
		ctrl      *gomock.Controller
		sut       flusher.Flusher
		input     []models.Remind
		chunkSize int
	)

	BeforeEach(func() {
		now := time.Now()
		chunkSize = 2
		ctrl = gomock.NewController(GinkgoT())
		mr = mocks.NewMockRepo(ctrl)
		sut = flusher.NewFlusher(mr, chunkSize)
		input = []models.Remind{
			0: models.NewRemind(1, 2, now, "birthday"),
			1: models.NewRemind(2, 2, now, "party"),
			2: models.NewRemind(3, 2, now, "meeting"),
			3: models.NewRemind(4, 2, now, "relax"),
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Flushing reminds (Flush method)", func() {
		Context("Flush without repo errors ", func() {
			It("There should be no reminders that are not saved to the repository", func() {
				mr.EXPECT().
					Add(gomock.Any()).
					Return(nil).
					AnyTimes()
				gomega.Expect(sut.Flush(input)).Should(gomega.BeEmpty(),
					"without errors")
			})
		})
		Context("Flush fails with repo error on every remind", func() {
			It("Output of Flush should be equal to input", func() {
				mr.EXPECT().
					Add(gomock.Any()).
					Return(errors.New("error on add")).
					AnyTimes()
				gomega.Expect(sut.Flush(input)).To(gomega.Equal(input),
					"all errors")
			})
		})
		Context("Flush fails with repo error only on first chunk", func() {
			It("Output of Flush should be the two first reminds", func() {
				mr.EXPECT().
					Add(gomock.Any()).
					Return(errors.New("error on add")).
					Times(1)
				mr.EXPECT().
					Add(gomock.Any()).
					Return(nil).
					AnyTimes()
				gomega.Expect(sut.Flush(input)).To(gomega.Equal(input[0:2]),
					"last two reminds saved without error")
			})
		})
	})
})
