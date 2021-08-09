package flusher_test

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestFlusher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flusher Suite")
}

var GinkgoT = ginkgo.GinkgoT
var RunSpecs = ginkgo.RunSpecs
var Fail = ginkgo.Fail
var Describe = ginkgo.Describe
var Context = ginkgo.Context
var It = ginkgo.It
var BeforeEach = ginkgo.BeforeEach
var AfterEach = ginkgo.AfterEach

var RegisterFailHandler = gomega.RegisterFailHandler
