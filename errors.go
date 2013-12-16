package cq

import (
	"errors"
	"log"
	"os"
)

var (
	errNotConnected             = errors.New("not connected")
	errNotImplemented           = errors.New("not implemented")
	errTransactionsNotSupported = errors.New("transactions aren't supported by your Neo4j version")

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
