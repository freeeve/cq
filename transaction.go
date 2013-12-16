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

type transactionResponse struct {
	Commit string `json:"commit"`
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

type commitError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type commitResponse struct {
	Errors []commitError `json:"errors"`
}

func (c *conn) Begin() (driver.Tx, error) {
	if c.transactionURL == "" {
		return nil, errTransactionsNotSupported
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(cypherTransaction{})
	req, err := http.NewRequest("POST", c.transactionURL, &buf)
	if err != nil {
		return nil, err
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	transResponse := transactionResponse{}
	json.NewDecoder(res.Body).Decode(&transResponse)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	c.transaction = &cypherTransaction{
		commitURL:      transResponse.Commit,
		transactionURL: res.Header.Get("Location"),
		c:              c,
	}
	c.transactionState = transactionStarted
	return c.transaction, nil
}

func (tx *cypherTransaction) query(query *string, args []driver.Value) error {
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
			return err
		}
	}
	return nil
}

func (tx *cypherTransaction) exec() error {
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

func (tx *cypherTransaction) Commit() error {
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
