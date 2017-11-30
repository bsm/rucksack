// Package log provides a 12-factor convenience wrapper around zap
package log

import (
	"fmt"
	"os"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var stdL atomic.Value

func init() {
	Replace(zap.NewExample())

	fields := parseFields(os.Getenv("LOG_FIELDS"))
	if fields == nil {
		fields = parseFields(os.Getenv("LOG_TAGS"))
	}

	logger, err := buildLogger(os.Getenv("LOG_NAME"), os.Getenv("LOG_LEVEL"), fields)
	if err == nil {
		Replace(logger)
	}
}

// --------------------------------------------------------------------

// L returns the global logger
func L() *zap.Logger {
	return stdL.Load().(*zap.Logger)
}

// Silence silences log output, useful in tests.
func Silence() { Replace(zap.NewNop()) }

// Replace allows you to replace the global logger. Example:
//
//   scoped := log.L().Named("worker").With(
//     zap.String("key", "value"),
//   )
//   log.Replace(scoped)
func Replace(l *zap.Logger) {
	stdL.Store(l)
}

// Standard logging methods

func Debug(args ...interface{}) { L().Debug(fmt.Sprint(args...)) }
func Error(args ...interface{}) { L().Error(fmt.Sprint(args...)) }
func Info(args ...interface{})  { L().Info(fmt.Sprint(args...)) }
func Warn(args ...interface{})  { L().Warn(fmt.Sprint(args...)) }
func Fatal(args ...interface{}) { L().Fatal(fmt.Sprint(args...)) }

func Debugf(format string, args ...interface{}) { L().Debug(fmt.Sprintf(format, args...)) }
func Errorf(format string, args ...interface{}) { L().Error(fmt.Sprintf(format, args...)) }
func Infof(format string, args ...interface{})  { L().Info(fmt.Sprintf(format, args...)) }
func Warnf(format string, args ...interface{})  { L().Warn(fmt.Sprintf(format, args...)) }
func Fatalf(format string, args ...interface{}) { L().Fatal(fmt.Sprintf(format, args...)) }

func Debugln(args ...interface{}) { L().Debug(fmt.Sprintln(args...)) }
func Errorln(args ...interface{}) { L().Error(fmt.Sprintln(args...)) }
func Infoln(args ...interface{})  { L().Info(fmt.Sprintln(args...)) }
func Warnln(args ...interface{})  { L().Warn(fmt.Sprintln(args...)) }
func Fatalln(args ...interface{}) { L().Fatal(fmt.Sprintln(args...)) }

func Debugw(msg string, fields ...zapcore.Field) { L().Debug(msg, fields...) }
func Errorw(msg string, fields ...zapcore.Field) { L().Error(msg, fields...) }
func Infow(msg string, fields ...zapcore.Field)  { L().Info(msg, fields...) }
func Warnw(msg string, fields ...zapcore.Field)  { L().Warn(msg, fields...) }
func Fatalw(msg string, fields ...zapcore.Field) { L().Fatal(msg, fields...) }

// Panic-catcher methods

// ErrorOnPanic logs error if func panics. Usage:
//
//   func someFunc() {
//     defer log.ErrorOnPanic()
//     ... // code that may panic()
//   }
func ErrorOnPanic() {
	if r := recover(); r != nil {
		Error("panic", zap.String("cause", fmt.Sprint(r)), zap.Stack("stack"))
	}
}

// FatalOnPanic calls Fatal if func panics. Usage:
//
//   func someFunc() {
//     defer log.FatalOnPanic()
//     ... // code that may panic()
//   }
func FatalOnPanic() {
	if r := recover(); r != nil {
		Fatal("panic", zap.String("cause", fmt.Sprint(r)), zap.Stack("stack"))
	}
}

// Sync performs a final sync. Best added to you main function. Usage:
//
//   func main() {
//     defer log.Sync()
//     ...
//   }
func Sync() error {
	return L().Sync()
}
