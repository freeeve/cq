package cq

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	query *string
}

func (stmt *cypherStmt) Close() error {
	stmt.query = nil
	return nil
}

func (stmt *cypherStmt) Exec(args []driver.Value) (driver.Result, error) {
	if stmt.c.transactionState == transactionStarted {
		//		errLog.Print("in transaction... queueing")
		stmt.c.transaction.query(stmt.query, args)
	} else {
		//		errLog.Print("not in transaction... querying: ", stmt.c)
		rows, err := stmt.Query(args)
		defer rows.Close()
		// TODO add counts and error support
		return nil, err
	}
	// never hit
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
	Query  *string                `json:"query"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Set("X-Stream", "true")
	req.Header.Set("User-Agent", cqVersion)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

func (stmt *cypherStmt) Query(args []driver.Value) (driver.Rows, error) {
	if stmt.c.transactionState == transactionStarted {
		return nil, errors.New("transactions only support Exec")
	} else {
		// this only happens outside of a transaction
		cyphReq := cypherRequest{
			Query: stmt.query,
		}
		if len(args) > 0 {
			cyphReq.Params = make(map[string]interface{})
		}
		for idx, e := range args {
			cyphReq.Params[strconv.Itoa(idx)] = e
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(cyphReq)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("POST", stmt.c.cypherURL, &buf)
		if err != nil {
			return nil, err
		}
		setDefaultHeaders(req)
		res, err := client.Do(req)
		defer res.Body.Close()
		if err != nil {
			return nil, err
		}
		cyphRes := cypherResult{}
		err = json.NewDecoder(res.Body).Decode(&cyphRes)
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
		if err != nil {
			return nil, err
		}
		if cyphRes.ErrorMessage != "" {
			return nil, errors.New("Cypher error: " + cyphRes.ErrorMessage)
		}
		return &rows{stmt, &cyphRes, 0}, nil
	}
	// never hits
	return nil, nil
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
	// TODO handle transaction
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
	Statement  *string                `json:"statement"`
	Parameters map[string]interface{} `json:"parameters"`
}

type cypherTransaction struct {
	Statements     []cypherTransactionStatement `json:"statements"`
	commitURL      string
	transactionURL string
	expiration     time.Time
	c              *conn
	rows           []*rows
}

func (tx *cypherTransaction) query(query *string, args []driver.Value) {
	//	errLog.Print("appending query", query)
	stmt := cypherTransactionStatement{
		Statement:  query,
		Parameters: make(map[string]interface{}, len(args)),
	}
	for idx, e := range args {
		stmt.Parameters[strconv.Itoa(idx)] = e
	}
	tx.Statements = append(tx.Statements, stmt)
	if len(tx.Statements) >= 100 {
		err := tx.exec()
		if err != nil {
			errLog.Print(err)
		}
	}
}

func (tx *cypherTransaction) exec() error {
	if tx.c.transactionState != transactionStarted {
		return errTransactionNotStarted
	}
	//jsontx, _ := json.Marshal(tx)
	//errLog.Print("executing a partial batch", string(jsontx))
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(tx)
	if err != nil {
		errLog.Print(err)
		return err
	}
	req, err := http.NewRequest("POST", tx.transactionURL, &buf)
	if err != nil {
		errLog.Print(err)
		return err
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
	commit := commitResponse{}
	json.NewDecoder(res.Body).Decode(&commit)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	tx.Statements = tx.Statements[:0]
	if len(commit.Errors) > 0 {
		return errors.New("exec errors: " + fmt.Sprintf("%q", commit))
	}
	if err != nil {
		errLog.Print(err)
		return err
	}
	return nil
}

type commitError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type commitResponse struct {
	Errors []commitError `json:"errors"`
}

func (tx *cypherTransaction) Commit() error {
	//	errLog.Print("committing transaction:", tx)
	if tx.c.transactionState != transactionStarted {
		return errTransactionNotStarted
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(tx)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", tx.commitURL, &buf)
	if err != nil {
		return err
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	commit := commitResponse{}
	json.NewDecoder(res.Body).Decode(&commit)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	if len(commit.Errors) > 0 {
		return errors.New("commit errors: " + fmt.Sprintf("%q", commit))
	}
	tx.c.transactionState = transactionNotStarted
	tx.c.transaction = nil
	return nil
}

func (tx *cypherTransaction) Rollback() error {
	if tx.c.transactionState != transactionStarted {
		return errTransactionNotStarted
	}
	req, err := http.NewRequest("DELETE", tx.transactionURL, nil)
	if err != nil {
		return err
	}
	setDefaultHeaders(req)
	//res, err := client.Do(req)
	if err != nil {
		return err
	}

	tx.c.transactionState = transactionNotStarted
	tx.c.transaction = nil
	return nil
}
