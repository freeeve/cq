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
	baseURL          string
	cypherURL        string
	transactionURL   string
	transaction      *cypherTransaction
	transactionState int
}

type neo4jBase struct {
	Data string `json:"data"`
}

type neo4jData struct {
	Cypher      string `json:"cypher"`
	Transaction string `json:"transaction"`
	Version     string `json:"neo4j_version"`
}

func setDefaultHeaders(req *http.Request) {
	req.Header.Set("X-Stream", "true")
	req.Header.Set("User-Agent", cqVersion)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
}

func Open(baseURL string) (driver.Conn, error) {
	// TODO
	// cache the results of this lookup
	// add support for multiple hosts (cluster)
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

	c := &conn{}
	c.cypherURL = neoData.Cypher
	c.transactionURL = neoData.Transaction

	return c, nil
}

func (c *conn) Close() error {
	// TODO check if in transaction and rollback
	return nil
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	if c.cypherURL == "" {
		return nil, errNotConnected
	}

	stmt := &cypherStmt{
		c:     c,
		query: &query,
	}

	return stmt, nil
}
