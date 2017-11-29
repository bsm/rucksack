package gcloud

import (
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Formatter", func() {
	var subject *Formatter

	It("should format", func() {
		data, err := subject.Format(&logrus.Entry{
			Data:    logrus.Fields{"str": "value", "num": 34, "err": errors.New("doh!")},
			Time:    time.Unix(1515151515, 0),
			Level:   logrus.WarnLevel,
			Message: "something happened",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(MatchJSON(`{
      "timestamp": "2018-01-05T11:25:15Z",
      "severity": "WARNING",
      "textPayload": "something happened",
      "labels": {
        "err": "doh!",
        "num": 34,
        "str": "value"
      }
    }`))
		Expect(data).To(HaveSuffix("}\n"))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rucksack/log/gcloud")
}
