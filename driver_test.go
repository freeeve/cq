package cq_test

import (
	"database/sql"
	_ "github.com/wfreeman/cq"
	. "launchpad.net/gocheck"
	"log"
)

// This file is meant to hold integration tests where cq must be imported as _

type DriverSuite struct{}

var _ = Suite(&DriverSuite{})

func testConn() *sql.DB {
	db, err := sql.Open("neo4j-cypher", "http://localhost:7474/")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func prepareTest(query string) *sql.Stmt {
	db := testConn()
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

func prepareAndQuery(query string) *sql.Rows {
	stmt := prepareTest(query)
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func (s *DriverSuite) TestDbQuery(c *C) {
	db := testConn()
	rows, err := db.Query("return 1")
	c.Assert(err, IsNil)

	if rows == nil {
		c.Fatal("rows shouldn't be nil")
	}
}

func (s *DriverSuite) TestDbExec(c *C) {
	db := testConn()
	result, err := db.Exec("return 1")
	c.Assert(err, IsNil)

	if result == nil {
		c.Fatal("result should not be nil")
	}
}

func (s *DriverSuite) TestQuerySimple(c *C) {
	rows := prepareAndQuery("return 1")
	hasNext := rows.Next()
	if !hasNext {
		c.Fatal("no next!")
	}

	var test int
	err := rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != 1 {
		c.Fatal("test != 1")
	}
}

func (s *DriverSuite) TestQuerySimpleFloat(c *C) {
	rows := prepareAndQuery("return 1.2")
	rows.Next()
	var test float64
	err := rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != 1.2 {
		c.Fatal("test != 1.2")
	}
}

func (s *DriverSuite) TestQuerySimpleString(c *C) {
	rows := prepareAndQuery("return '123'")
	rows.Next()
	var test string
	err := rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != "123" {
		c.Fatal("test != '123'")
	}
}

func (s *DriverSuite) TestQueryIntParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(123)
	c.Assert(err, IsNil)

	rows.Next()
	var test int
	err = rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != 123 {
		c.Fatal("test != 123;", test)
	}
}

func (s *DriverSuite) TestQuerySimpleIntArray(c *C) {
	c.Skip("can't convert to arrays yet")
	rows := prepareAndQuery("return [1,2,3]")
	rows.Next()
	var test []int
	err := rows.Scan(&test)
	c.Assert(err, IsNil)

	if test[0] != 1 || test[1] != 2 || test[2] != 3 {
		c.Fatal("test != [1,2,3];", test)
	}
}
