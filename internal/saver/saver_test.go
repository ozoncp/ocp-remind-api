package saver_test

import (
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/ozoncp/ocp-remind-api/internal/mocks"
	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/ozoncp/ocp-remind-api/internal/saver"
)

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
			2: models.NewRemind(3, 2, now, "meeting"),
			3: models.NewRemind(4, 2, now, "relax"),
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Saving reminds", func() {
		Context("Save 2 reminds into 2 capacity saver", func() {
			It("Flusher should be called 1 time", func() {
				mockFlusher := mocks.NewMockFlusher(ctrl)
				sut = saver.NewSaver(2, mockFlusher, 1*time.Millisecond)
				gomega.Expect(sut.Save(input[0])).Should(BeNil())
				gomega.Expect(sut.Save(input[1])).Should(BeNil())
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).MinTimes(0).MaxTimes(1)
				sut.Close()
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).MinTimes(0).MaxTimes(0)

			})
		})
		Context("Save 2 reminds into 0 capacity saver", func() {
			It("Flusher should be called 0 times", func() {
				mockFlusher := mocks.NewMockFlusher(ctrl)
				sut = saver.NewSaver(0, mockFlusher, 1*time.Second)
				gomega.Expect(sut.Save(input[0])).ShouldNot(BeNil())
				gomega.Expect(sut.Save(input[1])).ShouldNot(BeNil())
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).Times(0)
				sut.Close()
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).Times(0)
			})
		})
		Context("Save 3 reminds into 2 capacity saver", func() {
			It("Flusher should be called 2 times", func() {
				mockFlusher := mocks.NewMockFlusher(ctrl)
				sut = saver.NewSaver(5, mockFlusher, 1*time.Millisecond)
				gomega.Expect(sut.Save(input[0])).Should(BeNil())
				gomega.Expect(sut.Save(input[1])).Should(BeNil())
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).MinTimes(0).MaxTimes(1)
				gomega.Expect(sut.Save(input[2])).Should(BeNil())
				sut.Close()
				mockFlusher.EXPECT().Flush(gomock.Any()).Return([]models.Remind{}).MinTimes(1).MaxTimes(2)
			})
		})
	})
})
