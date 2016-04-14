package log

import (
	"bufio"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
)

func init() {
	base := logrus.New()

	// Set level
	base.Level = logrus.InfoLevel
	if level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		base.Level = level
	}

	// Create main entry
	entry := base.WithFields(nil)

	// Parse tags
	if s := os.Getenv("LOG_TAGS"); s != "" {
		for _, tag := range strings.Split(s, ",") {
			if parts := strings.SplitN(tag, ":", 2); len(parts) == 2 && parts[0] != "" {
				entry = entry.WithField(parts[0], parts[1])
			}
		}
	}

	std = &logger{Entry: entry, base: base}
	runtime.SetFinalizer(std, func(l *logger) { _ = l.Close() })
}

type logger struct {
	*logrus.Entry

	name string
	base *logrus.Logger

	writers []io.WriteCloser
	hooks   []Hook
}

// Close closes the log with its writers and its hooks
func (l *logger) Close() (err error) {
	for _, writer := range l.writers {
		if e := writer.Close(); e != nil {
			err = e
		}
	}
	l.writers = l.writers[:0]

	for _, hook := range l.hooks {
		if e := hook.Close(); e != nil {
			err = e
		}
	}
	l.hooks = l.hooks[:0]

	return
}

func (l *logger) plainLogger(fields logrus.Fields, level logrus.Level) *log.Logger {
	out := newWriter(l.WithFields(fields), level)
	l.writers = append(l.writers, out)
	return log.New(out, "", 0)
}

func newWriter(entry *logrus.Entry, level logrus.Level) *io.PipeWriter {
	reader, writer := io.Pipe()
	callback := entry.Info
	switch level {
	case logrus.WarnLevel:
		callback = entry.Warn
	case logrus.ErrorLevel:
		callback = entry.Error
	}
	go func() {
		defer reader.Close()

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			callback(scanner.Text())
		}
	}()

	runtime.SetFinalizer(writer, func(w *io.PipeWriter) { _ = w.Close() })
	return writer
}
