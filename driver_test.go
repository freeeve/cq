package cq_test

import (
	"database/sql"
	_ "github.com/wfreeman/cq"
	"log"
	"testing"
)

func testConn() *sql.DB {
	db, err := sql.Open("neo4j-cypher", "http://localhost:7474")
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

func failIfErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuerySimple(t *testing.T) {
	rows := prepareAndQuery("return 1")
	hasNext := rows.Next()
	if !hasNext {
		t.Fatal("no next!")
	}

	var test int
	err := rows.Scan(&test)
	failIfErr(err, t)

	if test != 1 {
		t.Fatal("test != 1")
	}
}

func TestQuerySimpleFloat(t *testing.T) {
	rows := prepareAndQuery("return 1.2")
	rows.Next()
	var test float64
	err := rows.Scan(&test)
	failIfErr(err, t)

	if test != 1.2 {
		t.Fatal("test != 1.2")
	}
}

func TestQuerySimpleString(t *testing.T) {
	rows := prepareAndQuery("return '123'")
	rows.Next()
	var test string
	err := rows.Scan(&test)
	failIfErr(err, t)

	if test != "123" {
		t.Fatal("test != '123'")
	}
}

// TODO array conversion
/*
func TestQuerySimpleIntArray(t *testing.T) {
	rows := prepareAndQuery("return [1,2,3]")
	rows.Next()
	var test []int
	err := rows.Scan(&test)
	failIfErr(err, t)

	if test[0] != 1 || test[1] != 2 || test[2] != 3 {
		t.Fatal("test != [1,2,3];", test)
	}
} */
