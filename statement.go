package cq

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	transactionNotStarted = iota
	transactionStarted    = iota
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
	stmt.Query(args)
	// TODO add counts and error support
	return nil, nil
}

func (stmt *cypherStmt) NumInput() int {
	// TODO maybe parse query to give a real number
	return -1 // avoid sanity check
}

type cypherResult struct {
	Columns         []string        `json:"columns"`
	Data            [][]interface{} `json:"data"`
	ErrorMessage    string          `json:"message"`
	ErrorException  string          `json:"exception"`
	ErrorFullname   string          `json:"fullname"`
	ErrorStacktrace []string        `json:"stacktrace"`
}

type cypherRequest struct {
	Query  string                 `json:"query"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Set("X-Stream", "true")
	req.Header.Set("User-Agent", cqVersion)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

func (stmt *cypherStmt) Query(args []driver.Value) (driver.Rows, error) {
	// TODO check if we're in a transaction and use it
	cyphReq := cypherRequest{
		Query: stmt.query,
	}
	if len(args) > 0 {
		cyphReq.Params = make(map[string]interface{})
	}
	for idx, e := range args {
		cyphReq.Params[strconv.Itoa(idx)] = e
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
	setDefaultHeaders(req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	cyphRes := cypherResult{}
	err = json.NewDecoder(res.Body).Decode(&cyphRes)
	if err != nil {
		return nil, err
	}
	if cyphRes.ErrorMessage != "" {
		return nil, errors.New("Cypher error: " + cyphRes.ErrorMessage)
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

type cypherTransactionStatement struct {
	Statement  string `json:"statement"`
	Parameters map[string]interface{}
}

type cypherTransaction struct {
	Statements     []cypherTransactionStatement `json:"statements"`
	commitURL      string
	transactionURL string
	expiration     time.Time
	c              *conn
}

func (tx *cypherTransaction) query(query string, args []driver.Value) {
	stmt := cypherTransactionStatement{
		Statement: query,
		//	Parameters: args,
	}
	tx.Statements = append(tx.Statements, stmt)
}

func (tx *cypherTransaction) Commit() error {
	if tx.c.transactionState != transactionStarted {
		return errTransactionNotStarted
	}
	// TODO commit

	tx.c.transactionState = transactionNotStarted
	tx.c.transaction = nil
	tx.c = nil
	return nil
}

func (tx *cypherTransaction) Rollback() error {
	if tx.c.transactionState != transactionStarted {
		return errTransactionNotStarted
	}
	// TODO rollback

	tx.c.transactionState = transactionNotStarted
	tx.c.transaction = nil
	tx.c = nil
	return nil
}
