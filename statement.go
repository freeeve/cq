package cq

import (
	"database/sql/driver"
)

type cypherStmt struct {
	c          *conn
	query      string
	paramCount int
}

func (stmt *cypherStmt) Close() error {
	return nil
}

func (stmt *cypherStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, nil
}

func (stmt *cypherStmt) NumInput() int {
	return 0
}

func (stmt *cypherStmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, nil
}
