package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type RunnerLogger struct {
	LogrusLogger *logrus.Logger
}

var (
	_std = RunnerLogger{
		LogrusLogger: logrus.StandardLogger(),
	}
	_mx sync.RWMutex
)

func getLogger() *logrus.Logger {
	_mx.RLock()
	defer _mx.RUnlock()
	return _std.LogrusLogger
}

func ChangeLogger(l RunnerLogger) {
	_mx.Lock()
	defer _mx.Unlock()
	_std = l
}

func CreateNewLogger() RunnerLogger {
	l := logrus.New()
	//l.SetLevel(defaultLevel)
	return RunnerLogger{l}
}

func WithError(err error) *logrus.Entry {
	return getLogger().WithError(err)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return getLogger().WithField(key, value)
}
func WithContext(ctx context.Context) *logrus.Entry {
	return getLogger().WithContext(ctx)
}
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return getLogger().WithFields(fields)
}
func WithTime(t time.Time) *logrus.Entry { return getLogger().WithTime(t) }

func Trace(args ...interface{}) { getLogger().Trace(args...) }
func Debug(args ...interface{}) { getLogger().Debug(args...) }
func Print(args ...interface{}) { getLogger().Print(args...) }
func Info(args ...interface{})  { getLogger().Info(args...) }
func Warn(args ...interface{})  { getLogger().Warn(args...) }
func Error(args ...interface{}) { getLogger().Error(args...) }
func Panic(args ...interface{}) { getLogger().Panic(args...) }
func Fatal(args ...interface{}) { getLogger().Fatal(args...) }

func Tracef(format string, args ...interface{}) {
	getLogger().Tracef(format, args...)
}
func Debugf(format string, args ...interface{}) {
	getLogger().Debugf(format, args...)
}
func Printf(format string, args ...interface{}) {
	getLogger().Printf(format, args...)
}
func Infof(format string, args ...interface{}) { getLogger().Infof(format, args...) }
func Warnf(format string, args ...interface{}) { getLogger().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) {
	getLogger().Errorf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	getLogger().Panicf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	getLogger().Fatalf(format, args...)
}

func Traceln(args ...interface{}) { getLogger().Traceln(args...) }
func Debugln(args ...interface{}) { getLogger().Debugln(args...) }
func Println(args ...interface{}) { getLogger().Println(args...) }
func Infoln(args ...interface{})  { getLogger().Infoln(args...) }
func Warnln(args ...interface{})  { getLogger().Warnln(args...) }
func Errorln(args ...interface{}) { getLogger().Errorln(args...) }
func Panicln(args ...interface{}) { getLogger().Panicln(args...) }
func Fatalln(args ...interface{}) { getLogger().Fatalln(args...) }

type contextKey string

const loggerKey contextKey = "global-logger"

// WithLogger adds a logrus logger to the context
func WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, _std)
}

// FromContext retrieves the logrus logger from the context
func FromContext(ctx context.Context) *logrus.Logger {
	_, ok := ctx.Value(loggerKey).(RunnerLogger)
	if !ok {
		// Return a default logger if none is found
		return logrus.New()
	}
	return getLogger()
}
