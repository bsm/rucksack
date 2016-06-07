package kafka

import (
	"testing"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kafka", func() {

	It("should extract tags", func() {
		tags := buildTags(mockEnv(map[string]string{
			"HOST":           "host.aws.dc",
			"PORT":           "8080",
			"LOG_KAFKA_TAGS": "foo:bar,oth:baz",
		}), 2233)
		Expect(tags).To(Equal(map[string]string{
			"host": "host",
			"port": "8080",
			"pid":  "2233",
			"foo":  "bar",
			"oth":  "baz",
		}))

		tags = buildTags(mockEnv(map[string]string{
			"HOST":           "host.aws.dc",
			"PORT":           "8080",
			"LOG_KAFKA_TAGS": "port:-,pid:",
		}), 2233)
		Expect(tags).To(Equal(map[string]string{
			"host": "host",
		}))
	})

	It("should configure levels", func() {
		f := &logrus.JSONFormatter{}
		p := newProducer(nil, "topic", logrus.WarnLevel, f, nil)
		Expect(p.Levels()).To(Equal([]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		}))
	})

})

// --------------------------------------------------------------------

func mockEnv(kv map[string]string) func(string) string {
	return func(key string) string {
		return kv[key]
	}
}

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rucksack/log/kafka")
}
