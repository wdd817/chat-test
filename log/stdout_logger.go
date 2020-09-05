package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type stdOutLogger struct {
	log.Logger
	LogLevel
}

func NewStdOutLogger() Logger {
	logger := &stdOutLogger{}
	logger.SetOutput(os.Stdout)
	logger.SetFlags(log.LstdFlags)
	return logger
}

func (logger *stdOutLogger) SetLevel(level LogLevel) {
	logger.LogLevel = level
}

// func (logger *stdOutLogger) Debug(args ...interface{}) {
// 	logger.print(DebugLevel, args)
// }
//
// func (logger *stdOutLogger) Info(args ...interface{}) {
// 	logger.print(InfoLevel, args)
// }
//
// func (logger *stdOutLogger) Warn(args ...interface{}) {
// 	logger.print(WarnLevel, args)
// }
//
// func (logger *stdOutLogger) Error(args ...interface{}) {
// 	logger.print(ErrorLevel, args)
// }
//
// func (logger *stdOutLogger) Fatal(args ...interface{}) {
// 	logger.print(FatalLevel, args)
// }

func (logger *stdOutLogger) Info(format string, args ...interface{}) {
	logger.printf(InfoLevel, format, args...)
}

func (logger *stdOutLogger) Debug(format string, args ...interface{}) {
	logger.printf(DebugLevel, format, args...)
}

func (logger *stdOutLogger) Warn(format string, args ...interface{}) {
	logger.printf(WarnLevel, format, args...)
}

func (logger *stdOutLogger) Error(format string, args ...interface{}) {
	logger.printf(ErrorLevel, format, args...)
}

func (logger *stdOutLogger) Fatal(format string, args ...interface{}) {
	logger.printf(FatalLevel, format, args...)
	os.Exit(-1)
}

func (logger *stdOutLogger) printf(level LogLevel, format string, args ...interface{}) {
	if level < logger.LogLevel {
		return
	}

	format = levelPrefix[level] + format
	logger.Logger.Output(4, fmt.Sprintf(format, args...)+logger.callFileLine())
}

func (logger *stdOutLogger) print(level LogLevel, args ...interface{}) {
	if level < logger.LogLevel {
		return
	}

	logger.Logger.Output(4, levelPrefix[level]+fmt.Sprint(args...)+logger.callFileLine())
}

func (logger *stdOutLogger) callFileLine() string {
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "???"
		line = 0
	}

	return fmt.Sprintf(" [%s:%d]", file, line)
}
