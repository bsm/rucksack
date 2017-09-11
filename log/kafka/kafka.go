/*
Package kafka forwards logs to kafka, if included.

  package main

  import(
    "github.com/bsm/rucksack/log"
    _ "github.com/bsm/rucksack/log/kafka"
  )

  func main() {
    sleep := make(chan struct{})
    <-sleep
  }

Run with:

  LOG_KAFKA_TOPIC=myapp LOG_KAFKA_ADDRS=broker-1:9092,broker-2:9092 go run main.go

Other env vars:

  LOG_KAFKA_LEVEL - configures the minimum level required for messages to be sent to kafka.
  Possible values are: PANIC, FATAL, ERROR, WARN, INFO, DEBUG. Default: INFO.

  LOG_KAFKA_TAGS - comma-separated key-value pairs, e.g. "foo:bar,oth:baz". By setting a value
  to "-", you can also suppress default tags, e.g. "port:-,host:-,pid:-".

*/
package kafka

import (
	"bytes"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
	"github.com/bsm/rucksack/log"
)

type producer struct {
	sarama.AsyncProducer
	topic  string
	tags   map[string]string
	format logrus.Formatter
	levels []logrus.Level
}

func NewProducer(addrs []string, topic string, level logrus.Level, format logrus.Formatter, tags map[string]string) (log.Hook, error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.NoResponse
	conf.Producer.Compression = sarama.CompressionSnappy
	conf.Producer.Return.Successes = false
	conf.Producer.Return.Errors = false
	conf.ChannelBufferSize = 1024

	ap, err := sarama.NewAsyncProducer(addrs, conf)
	if err != nil {
		return nil, err
	}
	return newProducer(ap, topic, level, format, tags), nil
}

func newProducer(ap sarama.AsyncProducer, topic string, level logrus.Level, format logrus.Formatter, tags map[string]string) log.Hook {
	var levels []logrus.Level
	for _, l := range logrus.AllLevels {
		if l <= level {
			levels = append(levels, l)
		}
	}

	return &producer{
		AsyncProducer: ap,

		format: format,
		topic:  topic,
		tags:   tags,
		levels: levels,
	}
}

// Levels implements github.com/bsm/rucksack/log.Hook
func (p *producer) Levels() []logrus.Level {
	return p.levels
}

// Fire implements github.com/bsm/rucksack/log.Hook
func (p *producer) Fire(entry *logrus.Entry) error {
	data := make(logrus.Fields, len(entry.Data)+len(p.tags))
	for k, v := range p.tags {
		data[k] = v
	}
	for k, v := range entry.Data {
		data[k] = v
	}

	line, err := p.format.Format(&logrus.Entry{
		Logger: entry.Logger,
		Data:   data,
	})
	if err != nil {
		return err
	}

	p.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(bytes.TrimSpace(line)),
	}
	return nil
}

// --------------------------------------------------------------------

func init() {
	addrs := os.Getenv("LOG_KAFKA_ADDRS")
	if addrs == "" {
		log.Error("log/kafka: LOG_KAFKA_ADDRS is required")
		return
	}

	topic := os.Getenv("LOG_KAFKA_TOPIC")
	if topic == "" {
		log.Error("log/kafka: LOG_KAFKA_TOPIC is required")
		return
	}

	tags := buildTags(os.Getenv, os.Getpid())
	level, err := logrus.ParseLevel("LOG_KAFKA_LEVEL")
	if err != nil {
		level = logrus.InfoLevel
	}

	format := &logrus.JSONFormatter{}
	hook, err := NewProducer(strings.Split(addrs, ","), topic, level, format, tags)
	if err != nil {
		log.Errorf("log/kafka: %s", err.Error())
		return
	}

	log.AddHook(hook)
}

func buildTags(env func(string) string, pid int) map[string]string {
	tags := make(map[string]string)
	if host := hostname(env("HOST")); host != "" {
		tags["host"] = host
	}
	if port := env("PORT"); port != "" {
		tags["port"] = port
	}
	if pid != 0 {
		tags["pid"] = strconv.Itoa(pid)
	}
	if extra := env("LOG_KAFKA_TAGS"); extra != "" {
		for _, pair := range strings.Split(extra, ",") {
			kv := strings.SplitN(pair, ":", 2)
			if kv[1] == "" || kv[1] == "-" {
				delete(tags, kv[0])
			} else if kv[0] != "" {
				tags[kv[0]] = kv[1]
			}
		}
	}
	return tags
}

func hostname(host string) string {
	if host == "" {
		host, _ = os.Hostname()
	}
	if pos := strings.Index(host, "."); pos > -1 {
		host = host[:pos]
	}
	return host
}
