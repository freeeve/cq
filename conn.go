// Package cq provides a database/sql implementation for Neo4j's Cypher query language.
package cq

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type cypherDriver struct{}

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
	baseURL        string
	cypherURL      string
	transactionURL string
	transaction    *cypherTransaction
}

type neo4jBase struct {
	Data string `json:"data"`
}

type neo4jData struct {
	Cypher      string `json:"cypher"`
	Transaction string `json:"transaction,omitempty"`
	Version     string `json:"neo4j_version"`
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Set("X-Stream", "true")
	req.Header.Set("User-Agent", cqVersion)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

// Open queries the base URL given to it for the Cypher
// and (optional) Transaction endpoints.
// It returns a connection handle, or an error if something went wrong.
func Open(baseURL string) (driver.Conn, error) {
	// TODO
	// cache the results of this lookup
	// add support for multiple hosts (cluster)
	neoBase, err := getNeoBase(baseURL)
	if err != nil {
		return nil, err
	}

	neoData, err := getNeoData(neoBase.Data)
	if err != nil {
		return nil, err
	}

	c := &conn{}
	c.cypherURL = neoData.Cypher
	c.transactionURL = neoData.Transaction

	return c, nil
}

func getNeoBase(url string) (*neo4jBase, error) {
	res, err := http.Get(url)
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
	return &neoBase, nil
}

func getNeoData(url string) (*neo4jData, error) {
	res, err := http.Get(url)
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
	return &neoData, nil
}

func (c *conn) Close() error {
	// TODO check if in transaction and rollback
	return nil
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	if c.cypherURL == "" {
		return nil, ErrNotConnected
	}

	stmt := &cypherStmt{
		c:     c,
		query: &query,
	}

	return stmt, nil
}
