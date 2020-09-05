package log

type LogLevel uint8

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelPrefix = []string{
	"[D] ",
	"[I] ",
	"[W] ",
	"[E] ",
}

type Logger interface {
	SetLevel(level LogLevel)

	// Debug(args ...interface{})
	// Info(args ...interface{})
	// Warn(args ...interface{})
	// Error(args ...interface{})
	// Fatal(args ...interface{})

	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

var gLogger Logger

func init() {
	gLogger = NewStdOutLogger()
}

// // Debug .
// func Debug(args ...interface{}) {
// 	gLogger.Debug(args)
// }
//
// // Info .
// func Info(args ...interface{}) {
// 	gLogger.Info(args)
// }
//
// // Warn .
// func Warn(args ...interface{}) {
// 	gLogger.Warn(args)
// }
//
// // Error .
// func Error(args ...interface{}) {
// 	gLogger.Error(args)
// }
//
// // Fatal .
// func Fatal(args ...interface{}) {
// 	gLogger.Fatal(args)
// }

// Debugf .
func Debug(format string, args ...interface{}) {
	gLogger.Debug(format, args...)
}

// Infof .
func Info(format string, args ...interface{}) {
	gLogger.Info(format, args...)
}

// Warnf .
func Warn(format string, args ...interface{}) {
	gLogger.Warn(format, args...)
}

// Errorf .
func Error(format string, args ...interface{}) {
	gLogger.Error(format, args...)
}

// Fatalf .
func Fatal(format string, args ...interface{}) {
	gLogger.Fatal(format, args...)
}
