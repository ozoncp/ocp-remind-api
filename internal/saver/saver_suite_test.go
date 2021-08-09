package saver_test

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestSaver(t *testing.T) { //nolint:paralleltest
	RegisterFailHandler(Fail)
	RunSpecs(t, "Saver Suite")
}

// Declarations for Ginkgo DSL

var GinkgoT = ginkgo.GinkgoT   //nolint:gofumpt,gochecknoglobals
var RunSpecs = ginkgo.RunSpecs //nolint:gochecknoglobals

var Fail = ginkgo.Fail //nolint:gochecknoglobals

var Describe = ginkgo.Describe //nolint:gochecknoglobals

var Context = ginkgo.Context //nolint:gochecknoglobals

var It = ginkgo.It //nolint:gochecknoglobals

var BeforeEach = ginkgo.BeforeEach //nolint:gochecknoglobals

var AfterEach = ginkgo.AfterEach //nolint:gochecknoglobals

// Declarations for Gomega DSL
var RegisterFailHandler = gomega.RegisterFailHandler

// Declarations for Gomega Matchers

var BeNil = gomega.BeNil
