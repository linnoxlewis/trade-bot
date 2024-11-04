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
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Llongfile),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.LstdFlags),
	}
}
