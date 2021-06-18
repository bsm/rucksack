package runtime

import (
	"testing"

	. "github.com/bsm/ginkgo"
	. "github.com/bsm/gomega"
)

var _ = Describe("runtime", func() {

	It("should extract levels", func() {
		Expect(levels).To(BeEmpty())
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rucksack/met/runtime")
}
