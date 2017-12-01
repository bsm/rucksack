package rucksack

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Tags",
	func(s string, exp []string) {
		Expect(Tags(s)).To(Equal(exp))
	},

	Entry("blank", "", []string{}),
	Entry("simple", "a,b", []string{"a", "b"}),
	Entry("spaced", "a ,  b", []string{"a", "b"}),
)

var _ = DescribeTable("Fields",
	func(s string, exp map[string]interface{}) {
		Expect(Fields(s)).To(Equal(exp))
	},

	Entry("blank", "", (map[string]interface{})(nil)),
	Entry("simple", "k1:v1,k2:v2", map[string]interface{}{"k1": "v1", "k2": "v2"}),
	Entry("with spaces", " k1:v1 ,   k2:v2", map[string]interface{}{"k1": "v1", "k2": "v2"}),
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rucksack")
}
