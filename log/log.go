// Package log provides a 12-factor convenience wrapper around logrus
package log

import (
	"io"
	"io/ioutil"
	"log"
	"runtime/debug"

	"github.com/Sirupsen/logrus"
)

// Hook is a mountable hook
type Hook interface {
	logrus.Hook
	Close() error
}

var std *logger

// Silence silences log output, useful in tests
func Silence() { SetOutput(ioutil.Discard) }

// NewInfo creates a stdlib *log.Logger, always logging at INFO level
func NewInfo(fields logrus.Fields) *log.Logger {
	return std.plainLogger(fields, logrus.InfoLevel)
}

// NewError creates a stdlib *log.Logger, always logging at ERROR level
func NewError(fields logrus.Fields) *log.Logger {
	return std.plainLogger(fields, logrus.ErrorLevel)
}

// SetOutput sets a custom output. Not thread-safe
func SetOutput(w io.Writer) { std.base.Out = w }

// SetFormatter sets a custom formatter. Not thread-safe
func SetFormatter(f logrus.Formatter) { std.base.Formatter = f }

// AddHook installs a custom hook to the logger. Not thread-safe
// TODO: use std.base.AddHook(hook) once https://github.com/sirupsen/logrus/pull/612 is released.
func AddHook(hook Hook) {
	std.hooks = append(std.hooks, hook)
	std.base.Hooks.Add(hook)
}

// Logging methods

func Debug(args ...interface{})                             { std.Debug(args...) }
func Debugf(format string, args ...interface{})             { std.Debugf(format, args...) }
func Debugln(args ...interface{})                           { std.Debugln(args...) }
func Error(args ...interface{})                             { std.Error(args...) }
func Errorf(format string, args ...interface{})             { std.Errorf(format, args...) }
func Errorln(args ...interface{})                           { std.Errorln(args...) }
func Info(args ...interface{})                              { std.Info(args...) }
func Infof(format string, args ...interface{})              { std.Infof(format, args...) }
func Infoln(args ...interface{})                            { std.Infoln(args...) }
func Print(args ...interface{})                             { std.Print(args...) }
func Printf(format string, args ...interface{})             { std.Printf(format, args...) }
func Println(args ...interface{})                           { std.Println(args...) }
func Fatal(args ...interface{})                             { std.Fatal(args...) }
func Fatalf(format string, args ...interface{})             { std.Fatalf(format, args...) }
func Fatalln(args ...interface{})                           { std.Fatalln(args...) }
func Warn(args ...interface{})                              { std.Warn(args...) }
func Warnf(format string, args ...interface{})              { std.Warnf(format, args...) }
func Warnln(args ...interface{})                            { std.Warnln(args...) }
func WithField(key string, value interface{}) *logrus.Entry { return std.WithField(key, value) }
func WithFields(fields logrus.Fields) *logrus.Entry         { return std.WithFields(fields) }

// Panic-catcher methods

// ErrorOnPanic logs error if func panics.
//
// Usage:
//   func someFunc() {
//     defer ErrorOnPanic()
//     ... // code that may panic()
//   }
func ErrorOnPanic() {
	if r := recover(); r != nil {
		Errorf("panic: %v\n\n%s", r, debug.Stack())
	}
}

// FatalOnPanic calls Fatal if func panics.
//
// Usage:
//   func someFunc() {
//     defer FatalOnPanic()
//     ... // code that may panic()
//   }
func FatalOnPanic() {
	if r := recover(); r != nil {
		Fatalf("panic: %v\n\n%s", r, debug.Stack())
	}
}
