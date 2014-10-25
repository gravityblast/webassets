package main

import (
	"log"
	"os"
)

type simpleLogger struct {
	verbose bool
	logger  *log.Logger
}

var logger *simpleLogger

func (l *simpleLogger) Printf(format string, values ...interface{}) {
	if l.verbose {
		l.logger.Printf(format, values...)
	}
}

func (l *simpleLogger) Fatalf(format string, values ...interface{}) {
	l.logger.Fatalf(format, values...)
}

func init() {
	logger = &simpleLogger{
		logger: log.New(os.Stderr, "", 0),
	}
}
