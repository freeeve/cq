package cq

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	transactionNotStarted = iota
	transactionStarted    = iota
)

type transactionResponse struct {
	Commit      string `json:"commit"`
	Transaction struct {
		Expires string
	}
	Errors []commitError `json:"errors"`
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
	keepAlive      *time.Timer
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
	exp, err := time.Parse(time.RFC1123Z, transResponse.Transaction.Expires)
	if err != nil {
		log.Println(err, c)
		err = nil
	}
	c.transaction = &cypherTransaction{
		commitURL:      transResponse.Commit,
		transactionURL: res.Header.Get("Location"),
		c:              c,
		expiration:     exp,
	}
	c.transaction.updateKeepAlive()
	return c.transaction, nil
}

func (tx *cypherTransaction) query(query *string, args []driver.Value) error {
	stmt := cypherTransactionStatement{
		Statement:  query,
		Parameters: makeArgsMap(args),
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
		return err
	}
	req, err := http.NewRequest("POST", tx.transactionURL, &buf)
	if err != nil {
		return err
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	trans := transactionResponse{}
	json.NewDecoder(res.Body).Decode(&trans)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()

	tx.expiration, err = time.Parse(time.RFC1123Z, trans.Transaction.Expires)
	if err != nil {
		log.Print(err, tx)
		err = nil
	}
	tx.updateKeepAlive()

	tx.Statements = tx.Statements[:0]

	if len(trans.Errors) > 0 {
		return errors.New("exec errors: " + fmt.Sprintf("%q", trans))
	}
	if err != nil {
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
	tx.c.transaction = nil
	if tx.keepAlive != nil {
		tx.keepAlive.Stop()
	}
	return nil
}

func (tx *cypherTransaction) Rollback() error {
	req, err := http.NewRequest("DELETE", tx.transactionURL, nil)
	if err != nil {
		return err
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
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
		return errors.New("rollback errors: " + fmt.Sprintf("%q", commit))
	}
	tx.c.transaction = nil
	if tx.keepAlive != nil {
		tx.keepAlive.Stop()
	}

	return nil
}

func (tx *cypherTransaction) updateKeepAlive() {
	if tx.keepAlive != nil {
		tx.keepAlive.Stop()
	}
	dur := -1 * time.Since(tx.expiration)
	if dur <= 1*time.Second {
		dur = 500 * time.Millisecond
	}
	tx.keepAlive = time.AfterFunc(dur, func() { sendKeepAlive(tx.transactionURL) })
}

func sendKeepAlive(txURL string) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(cypherTransaction{
		Statements: []cypherTransactionStatement{},
	})
	if err != nil {
		log.Print(err)
	}
	req, err := http.NewRequest("POST", txURL, &buf)
	if err != nil {
		log.Print(err)
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
	if err != nil {
	}
	trans := transactionResponse{}
	json.NewDecoder(res.Body).Decode(&trans)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if len(trans.Errors) > 0 {
		//log.Println(trans)
		return
	}
	exp, err := time.Parse(time.RFC1123Z, trans.Transaction.Expires)
	if err != nil {
		//log.Println(err)
		return
	}
	dur := -1 * time.Since(exp)
	time.AfterFunc(dur, func() { sendKeepAlive(txURL) })
}
