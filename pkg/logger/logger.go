package logger

import (
	"log"
)

type Logger struct{}

func New() *Logger { return &Logger{} }

func (l *Logger) Infof(format string, args ...any)  { log.Printf("INFO: "+format, args...) }
func (l *Logger) Warnf(format string, args ...any)  { log.Printf("WARN: "+format, args...) }
func (l *Logger) Errorf(format string, args ...any) { log.Printf("ERROR: "+format, args...) }
func (l *Logger) Fatalf(format string, args ...any) { log.Fatalf("FATAL: "+format, args...) }
