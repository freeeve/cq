package cq

import (
	"errors"
)

var (
	errNotConnected             = errors.New("cq: not connected")
	errNotImplemented           = errors.New("cq: not implemented")
	errTransactionsNotSupported = errors.New("cq: transactions aren't supported by your Neo4j version")
)
