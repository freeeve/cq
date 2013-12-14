package cq

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
)

type CypherDriver struct{}

var count int = 0

func (d *CypherDriver) Open(name string) (driver.Conn, error) {
	return Open(name)
}

func init() {
	sql.Register("neo4j-cypher", &CypherDriver{})
}

var (
	cqVersion = "0.1.0"
)

type conn struct {
	baseURL          string
	cypherURL        string
	transactionURL   string
	transaction      *cypherTransaction // for now going to support one transaction per connection
	transactionState int
	id               int
}

type Neo4jBase struct {
	Data string `json:"data"`
}

type Neo4jData struct {
	Cypher      string `json:"cypher"`
	Transaction string `json:"transaction"`
	Version     string `json:"neo4j_version"`
}

// TODO
// cache the results of this lookup
// add support for multiple hosts (cluster)
func Open(baseURL string) (driver.Conn, error) {
	resp, err := http.Get(baseURL)
   defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	neo4jBase := Neo4jBase{}
	err = json.NewDecoder(resp.Body).Decode(&neo4jBase)
	if err != nil {
		return nil, err
	}

   resp2, err := http.Get(neo4jBase.Data)
   defer resp2.Body.Close()
	if err != nil {
		return nil, err
	}

	neo4jData := Neo4jData{}
	err = json.NewDecoder(resp2.Body).Decode(&neo4jData)
	if err != nil {
		return nil, err
	}

	count++
	c := &conn{id: count}
	c.cypherURL = neo4jData.Cypher
	c.transactionURL = neo4jData.Transaction

	return c, nil
}

type TransactionResponse struct {
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
	client := &http.Client{}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(cypherTransaction{})
	req, err := http.NewRequest("POST", c.transactionURL, &buf)
	if err != nil {
		return nil, err
	}
	setDefaultHeaders(req)
	res, err := client.Do(req)
   defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	transactionResponse := TransactionResponse{}
	json.NewDecoder(res.Body).Decode(&transactionResponse)
	c.transaction = &cypherTransaction{
		commitURL:      transactionResponse.Commit,
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
		query: query,
	}

	return stmt, nil
}
