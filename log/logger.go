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

func Debug(format string, args ...interface{}) {
	gLogger.Info(format, args)
}

func Info(format string, args ...interface{}) {
	gLogger.Info(format, args)
}

func Warn(format string, args ...interface{}) {
	gLogger.Warn(format, args)
}

func Error(format string, args ...interface{}) {
	gLogger.Error(format, args)
}

func Fatal(format string, args ...interface{}) {
	gLogger.Fatal(format, args)
}
