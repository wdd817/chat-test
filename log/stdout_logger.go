package log

import (
	"fmt"
	"log"
	"os"
)

type stdOutLogger struct {
	log.Logger
	LogLevel
}

func NewStdOutLogger() Logger {
	logger := &stdOutLogger{}
	logger.SetOutput(os.Stdout)
	return logger
}

func (logger *stdOutLogger) SetLevel(level LogLevel) {
	logger.LogLevel = level
}

func (logger *stdOutLogger) Info(format string, args ...interface{}) {
	logger.print(InfoLevel, format, args)
}

func (logger *stdOutLogger) Debug(format string, args ...interface{}) {
	logger.print(DebugLevel, format, args)
}

func (logger *stdOutLogger) Warn(format string, args ...interface{}) {
	logger.print(WarnLevel, format, args)
}

func (logger *stdOutLogger) Error(format string, args ...interface{}) {
	logger.print(ErrorLevel, format, args)
}

func (logger *stdOutLogger) Fatal(format string, args ...interface{}) {
	logger.print(FatalLevel, format, args)
}

func (logger *stdOutLogger) print(level LogLevel, format string, args ...interface{}) {
	if level < logger.LogLevel {
		return
	}

	format = levelPrefix[level] + format
	_ = logger.Logger.Output(4, fmt.Sprint(format, args))
}
