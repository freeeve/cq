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

	errLog = log.New(os.Stderr, "[Cypher] ", log.Ldate|log.Ltime|log.Lshortfile)
)

// SetLogger allows users to set the logger used for errors from cq.
// It returns an error only when the logger is nil.
func SetLogger(logger *log.Logger) error {
	if logger == nil {
		return errors.New("logger is nil")
	}
	errLog = logger
	return nil
}
