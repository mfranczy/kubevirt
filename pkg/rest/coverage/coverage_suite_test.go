package coverage_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCoverage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coverage Suite")
}
