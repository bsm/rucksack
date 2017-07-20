package gcloud

import (
	"errors"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatter", func() {
	var subject *Formatter

	It("should format", func() {
		data, err := subject.Format(&logrus.Entry{
			Data:    logrus.Fields{"str": "value", "num": 34, "err": errors.New("doh!")},
			Time:    time.Unix(15151515, 0),
			Level:   logrus.WarnLevel,
			Message: "something happened",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(MatchJSON(`{
      "timestamp": "1970-06-25T09:45:15+01:00",
      "severity": "WARNING",
      "textPayload": "something happened",
      "labels": {
        "err": "doh!",
        "num": 34,
        "str": "value"
      }
    }`))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rucksack/log/gcloud")
}
