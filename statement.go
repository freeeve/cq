package cq

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"io"
	"net/http"
)

type cypherStmt struct {
	c     *conn
	query string
}

func (stmt *cypherStmt) Close() error {
	stmt.query = ""
	return nil
}

func (stmt *cypherStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errNotImplemented
}

func (stmt *cypherStmt) NumInput() int {
	// TODO maybe parse query to give a real number
	return -1 // avoid sanity check
}

type cypherResult struct {
	Columns []string        `json:"columns"`
	Data    [][]interface{} `json:"data"`
}

type cypherRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func (stmt *cypherStmt) Query(args []driver.Value) (driver.Rows, error) {
	cyphReq := cypherRequest{
		Query: stmt.query,
	}
	// TODO figure out how to use encoder for streaming here?
	buf, err := json.Marshal(cyphReq)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", stmt.c.cypherURL, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Stream", "true")
	req.Header.Set("User-Agent", cqVersion)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	cyphRes := cypherResult{}
	err = json.NewDecoder(res.Body).Decode(&cyphRes)
	if err != nil {
		return nil, err
	}
	return &rows{stmt, &cyphRes, 0}, nil
}

type rows struct {
	stmt   *cypherStmt
	result *cypherResult
	pos    int
}

func (rs *rows) Close() error {
	rs.result = nil
	return nil
}

func (rs *rows) Columns() []string {
	return rs.result.Columns
}

func (rs *rows) Next(dest []driver.Value) error {
	if len(rs.result.Data) <= rs.pos {
		return io.EOF
	}
	for i := 0; i < len(dest); i++ {
		dest[i] = rs.result.Data[rs.pos][i]
	}
	rs.pos++
	return nil
}
