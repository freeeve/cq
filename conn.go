package cq

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
)

type CypherDriver struct{}

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
	baseURL        string
	cypherURL      string
	transactionURL string
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
	if err != nil {
		return nil, err
	}

	neo4jBase := Neo4jBase{}
	err = json.NewDecoder(resp.Body).Decode(&neo4jBase)
	if err != nil {
		return nil, err
	}

	resp, err = http.Get(neo4jBase.Data)
	if err != nil {
		return nil, err
	}

	neo4jData := Neo4jData{}
	err = json.NewDecoder(resp.Body).Decode(&neo4jData)
	if err != nil {
		return nil, err
	}

	c := conn{}
	c.cypherURL = neo4jData.Cypher
	c.transactionURL = neo4jData.Transaction

	return c, nil
}

func (c conn) Begin() (driver.Tx, error) {
	return nil, errors.New("transaction mode not supported yet")
}

func (c conn) Close() error {
	// nothing needs done - connections aren't held open
	return nil
}

func (c conn) Prepare(query string) (driver.Stmt, error) {
	if c.cypherURL == "" {
		return nil, errNotConnected
	}

	stmt := &cypherStmt{
		c:     &c,
		query: query,
	}

	return stmt, nil
}
