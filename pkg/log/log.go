package log

import (
	"log"
	"os"
)

type Logger struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		ErrorLog: NewErrorLog(),
		InfoLog:  NewInfoLog(),
	}
}

func NewInfoLog() *log.Logger {
	return log.New(os.Stdout, "INFO\t", log.LstdFlags)
}

func NewErrorLog() *log.Logger {
	return log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Llongfile)
}
