package cq

import (
	"errors"
)

var (
	ErrNotConnected             = errors.New("cq: not connected")
	ErrNotImplemented           = errors.New("cq: not implemented")
	ErrTransactionsNotSupported = errors.New("cq: transactions aren't supported by your Neo4j version")
	ErrScanOnNil                = errors.New("cq: can't scan into nil pointer")
)
