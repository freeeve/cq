package cq

import (
	"errors"
	"log"
	"os"
)

var (
	errNotConnected   = errors.New("Not Connected")
	errNotImplemented = errors.New("Not Implemented")

	errLog Logger = log.New(os.Stderr, "[Cypher] ", log.Ldate|log.Ltime|log.Lshortfile)
)

type Logger interface {
	Print(v ...interface{})
}

func SetLogger(logger Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	errLog = logger
	return nil
}
