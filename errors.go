package cq

import (
	"errors"
)

var (
	errNotConnected             = errors.New("not connected")
	errNotImplemented           = errors.New("not implemented")
	errTransactionsNotSupported = errors.New("transactions aren't supported by your Neo4j version")
)
