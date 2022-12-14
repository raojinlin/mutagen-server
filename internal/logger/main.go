package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelDebug
	LevelError
)

type Logger struct {
	Level     Level
	Output    *os.File
	Name      string
	DebugMode bool
	loggers   map[Level]*log.Logger
}

func (l *Logger) init() {
	flag := log.Ldate | log.Ltime
	if l.DebugMode {
		flag |= log.Lshortfile
	}

	l.loggers[LevelInfo] = log.New(l.Output, l.getPrefix("INFO"), flag)
	l.loggers[LevelWarning] = log.New(l.Output, l.getPrefix("WARN"), flag)
	l.loggers[LevelDebug] = log.New(l.Output, l.getPrefix("DEBUG"), flag)
	l.loggers[LevelError] = log.New(l.Output, l.getPrefix("ERROR"), flag)
}

func (l *Logger) getPrefix(level string) string {
	return fmt.Sprintf("[%s] %s: ", l.Name, level)
}

func (l *Logger) Record(v ...interface{}) {
	output := map[string]interface{}{
		"Name":  l.Name,
		"Level": l.Level,
		"Debug": l.DebugMode,
		"Item":  v,
	}

	result, _ := json.Marshal(output)
	l.loggers[l.Level].Println(string(result))
}

func (l *Logger) Println(v ...interface{}) {
	l.loggers[l.Level].Println(v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.loggers[l.Level].Print(v...)
}

func (l *Logger) Printf(fmt string, v ...interface{}) {
	l.loggers[l.Level].Printf(fmt, v...)
}

func (l *Logger) Info() *log.Logger {
	return l.loggers[LevelInfo]
}

func (l *Logger) Error() *log.Logger {
	return l.loggers[LevelError]
}

func (l *Logger) Debug() *log.Logger {
	return l.loggers[LevelDebug]
}

func (l *Logger) Warn() *log.Logger {
	return l.loggers[LevelWarning]
}

func (l *Logger) Warning() *log.Logger {
	return l.Warn()
}

// NewLogger - create and init logger, if output is nil the print to os.Stdout
func NewLogger(name string, output *os.File, level Level, debug bool) *Logger {
	if output == nil {
		output = os.Stdout
	}

	logger := &Logger{
		Name:      name,
		Level:     level,
		Output:    output,
		DebugMode: debug,
		loggers:   make(map[Level]*log.Logger),
	}

	logger.init()
	return logger
}
