package cq

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type cypherDriver struct{}

var count int = 0

func (d *cypherDriver) Open(name string) (driver.Conn, error) {
	return Open(name)
}

func init() {
	sql.Register("neo4j-cypher", &cypherDriver{})
}

var (
	cqVersion = "0.1.0"
	tr        = &http.Transport{
		DisableKeepAlives: true,
	}
	client = &http.Client{}
)

type conn struct {
	baseURL          string
	cypherURL        string
	transactionURL   string
	transaction      *cypherTransaction // for now going to support one transaction per connection
	transactionState int
	id               int
}

type neo4jBase struct {
	Data string `json:"data"`
}

type neo4jData struct {
	Cypher      string `json:"cypher"`
	Transaction string `json:"transaction"`
	Version     string `json:"neo4j_version"`
}

// TODO
// cache the results of this lookup
// add support for multiple hosts (cluster)
func Open(baseURL string) (driver.Conn, error) {
	res, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}

	neoBase := neo4jBase{}
	err = json.NewDecoder(res.Body).Decode(&neoBase)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	res, err = http.Get(neoBase.Data)
	if err != nil {
		return nil, err
	}

	neoData := neo4jData{}
	err = json.NewDecoder(res.Body).Decode(&neoData)
	io.Copy(ioutil.Discard, res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	count++
	c := &conn{id: count}
	c.cypherURL = neoData.Cypher
	c.transactionURL = neoData.Transaction

	return c, nil
}

type transactionResponse struct {
	Commit string `json:"commit"`
}

func (c *conn) Begin() (driver.Tx, error) {
	if c.transactionURL == "" {
		return nil, errTransactionsNotSupported
	}
	if c.transactionState == transactionStarted {
		// this should not happen. probably delete this check (since a new connection will be allocated)
		return nil, errTransactionStarted
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
	//	errLog.Print("transaction successfully started:", c, c.transaction)
	return c.transaction, nil
}

func (c *conn) Close() error {
	// TODO check if in transaction and rollback
	return nil
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	//	errLog.Print("preparing a query: ", c)
	if c.cypherURL == "" {
		return nil, errNotConnected
	}

	stmt := &cypherStmt{
		c:     c,
		query: &query,
	}

	return stmt, nil
}
