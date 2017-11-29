// Package log provides a 12-factor convenience wrapper around zap
package log

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Hook is a mountable hook
type Hook interface {
	OnEntry(zapcore.Entry) error
	Close() error
}

var std = zap.NewNop()

func init() {
	name := os.Getenv("LOG_NAME")
	level := os.Getenv("LOG_LEVEL")
	fields := parseFields(os.Getenv("LOG_FIELDS"))
	if fields == nil {
		fields = parseFields(os.Getenv("LOG_TAGS"))
	}

	if logger, err := buildLogger(name, level, fields); err == nil {
		std = logger
	}
	runtime.SetFinalizer(std, func(l *zap.Logger) { _ = l.Sync() })
}

// --------------------------------------------------------------------

// Silence silences log output, useful in tests
func Silence() { std = zap.NewNop() }

// NewStdLogAt creates a stdlib *log.Logger at given level
func NewStdLogAt(level zapcore.Level, fields ...zapcore.Field) (*log.Logger, error) {
	return zap.NewStdLogAt(std.With(fields...), level)
}

// AddHook installs a custom hook to the logger.
func AddHook(hook Hook) {
	std = std.WithOptions(zap.Hooks(hook.OnEntry))
}

// Logging methods

func Print(args ...interface{})                 { Info(args...) }
func Printf(format string, args ...interface{}) { Infof(format, args...) }
func Println(args ...interface{})               { Infoln(args...) }

func Debug(args ...interface{}) { std.Debug(fmt.Sprint(args...)) }
func Error(args ...interface{}) { std.Error(fmt.Sprint(args...)) }
func Info(args ...interface{})  { std.Info(fmt.Sprint(args...)) }
func Warn(args ...interface{})  { std.Warn(fmt.Sprint(args...)) }
func Fatal(args ...interface{}) { std.Fatal(fmt.Sprint(args...)) }

func Debugf(format string, args ...interface{}) { std.Debug(fmt.Sprintf(format, args...)) }
func Errorf(format string, args ...interface{}) { std.Error(fmt.Sprintf(format, args...)) }
func Infof(format string, args ...interface{})  { std.Info(fmt.Sprintf(format, args...)) }
func Warnf(format string, args ...interface{})  { std.Warn(fmt.Sprintf(format, args...)) }
func Fatalf(format string, args ...interface{}) { std.Fatal(fmt.Sprintf(format, args...)) }

func Debugln(args ...interface{}) { std.Debug(fmt.Sprintln(args...)) }
func Errorln(args ...interface{}) { std.Error(fmt.Sprintln(args...)) }
func Infoln(args ...interface{})  { std.Info(fmt.Sprintln(args...)) }
func Fatalln(args ...interface{}) { std.Fatal(fmt.Sprintln(args...)) }
func Warnln(args ...interface{})  { std.Warn(fmt.Sprintln(args...)) }

func Debugw(msg string, fields ...zapcore.Field) { std.Debug(msg, fields...) }
func Errorw(msg string, fields ...zapcore.Field) { std.Error(msg, fields...) }
func Infow(msg string, fields ...zapcore.Field)  { std.Info(msg, fields...) }
func Warnw(msg string, fields ...zapcore.Field)  { std.Warn(msg, fields...) }
func Fatalw(msg string, fields ...zapcore.Field) { std.Fatal(msg, fields...) }

func With(fields ...zapcore.Field) *zap.Logger { return std.With(fields...) }

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
		Error("panic", zap.String("cause", fmt.Sprint(r)), zap.Stack("stack"))
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
		Fatal("panic", zap.String("cause", fmt.Sprint(r)), zap.Stack("stack"))
	}
}
