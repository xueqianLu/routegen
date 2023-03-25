package log

import (
	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var (
	mlog = logrus.New()
)

type LogConfig struct {
	Save  uint   `json:"save"`
	Level string `json:"level"`
}

func InitLog() {
	mlog.Out = os.Stdout
	mlog.SetLevel(logrus.DebugLevel)
	mlog.Formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-2 15:04:05.000"}
	//localFilesystemLogger(mlog, logConfig.Path, logConfig.Save)
}

func logWriter(logPath string, level string, save uint) *rotatelogs.RotateLogs {
	logFullPath := path.Join(logPath, level)
	logwriter, err := rotatelogs.New(
		logFullPath+".%Y%m%d",
		rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationCount(save),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return logwriter
}

func localFilesystemLogger(log *logrus.Logger, logPath string, save uint) {
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: logWriter(logPath, "debug", save), // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  logWriter(logPath, "info", save),
		logrus.WarnLevel:  logWriter(logPath, "warn", save),
		logrus.ErrorLevel: logWriter(logPath, "error", save),
		logrus.FatalLevel: logWriter(logPath, "fatal", save),
		logrus.PanicLevel: logWriter(logPath, "panic", save),
	}, &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-2 15:04:05.000"})
	log.AddHook(lfHook)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func WithField(key string, value interface{}) *logrus.Entry {
	return mlog.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return mlog.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func WithError(err error) *logrus.Entry {
	return mlog.WithError(err)
}

// Add a context to the log entry.
func WithContext(ctx context.Context) *logrus.Entry {
	return mlog.WithContext(ctx)
}

// Overrides the time of the log entry.
func WithTime(t time.Time) *logrus.Entry {
	return mlog.WithTime(t)
}

func Logf(level logrus.Level, format string, args ...interface{}) {
	mlog.Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	mlog.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	mlog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	mlog.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	mlog.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	mlog.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	mlog.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	mlog.Errorf(format, args)
}

func Fatalf(format string, args ...interface{}) {
	mlog.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	mlog.Panicf(format, args...)
}

func Log(level logrus.Level, args ...interface{}) {
	mlog.Log(level, args...)
}

func LogFn(level logrus.Level, fn logrus.LogFunction) {
	mlog.LogFn(level, fn)
}

func Trace(args ...interface{}) {
	mlog.Trace(args...)
}

func Debug(args ...interface{}) {
	mlog.Debug(args...)
}

func Info(args ...interface{}) {
	mlog.Info(args...)
}

func Print(args ...interface{}) {
	mlog.Print(args...)
}

func Warn(args ...interface{}) {
	mlog.Warn(args...)
}

func Warning(args ...interface{}) {
	mlog.Warning(args...)
}

func Error(args ...interface{}) {
	mlog.Error(args...)
}

func Fatal(args ...interface{}) {
	mlog.Fatal(args...)
}

func Panic(args ...interface{}) {
	mlog.Panic(args...)
}

func TraceFn(fn logrus.LogFunction) {
	mlog.TraceFn(fn)
}

func DebugFn(fn logrus.LogFunction) {
	mlog.DebugFn(fn)
}

func InfoFn(fn logrus.LogFunction) {
	mlog.InfoFn(fn)
}

func PrintFn(fn logrus.LogFunction) {
	mlog.PrintFn(fn)
}

func WarnFn(fn logrus.LogFunction) {
	mlog.WarnFn(fn)
}

func WarningFn(fn logrus.LogFunction) {
	mlog.WarningFn(fn)
}

func ErrorFn(fn logrus.LogFunction) {
	mlog.ErrorFn(fn)
}

func FatalFn(fn logrus.LogFunction) {
	mlog.FatalFn(fn)
}

func PanicFn(fn logrus.LogFunction) {
	mlog.PanicFn(fn)
}

func Logln(level logrus.Level, args ...interface{}) {
	mlog.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	mlog.Traceln(args...)
}

func Debugln(args ...interface{}) {
	mlog.Debugln(args...)
}

func Infoln(args ...interface{}) {
	mlog.Infoln(args...)
}

func Println(args ...interface{}) {
	mlog.Println(args...)
}

func Warnln(args ...interface{}) {
	mlog.Warnln(args...)
}

func Warningln(args ...interface{}) {
	mlog.Warningln(args...)
}

func Errorln(args ...interface{}) {
	mlog.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	mlog.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	mlog.Panicln(args...)
}
