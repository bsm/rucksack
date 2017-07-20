/*
Package gcloud includes Google Cloud helpers

  package main

  import(
    "github.com/bsm/rucksack/log"
    "github.com/bsm/rucksack/log/gcloud"
  )

  func init() {
    log.SetFormatter(new(gcloud.Formatter))
  }

*/
package gcloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

var _ logrus.Formatter = (*Formatter)(nil)

// Formatter is a Google Cloud Logging compatible formatter
type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := make(logrus.Fields, len(entry.Data))
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			fields[k] = v.Error()
		default:
			fields[k] = v
		}
	}

	line, err := json.Marshal(&loggingEntry{
		Timestamp: entry.Time,
		Severity:  strings.ToUpper(entry.Level.String()),
		Payload:   entry.Message,
		Labels:    fields,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(line, '\n'), nil
}

type loggingEntry struct {
	Timestamp time.Time     `json:"timestamp"`
	Severity  string        `json:"severity,omitempty"`
	Payload   string        `json:"textPayload,omitempty"`
	Labels    logrus.Fields `json:"labels,omitempty"`
}
