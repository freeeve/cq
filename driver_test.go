package cq_test

import (
	"database/sql"
	_ "github.com/wfreeman/cq"
	"log"
	"testing"
)

func testConn() *sql.DB {
	db, err := sql.Open("neo4j-cypher", "http://127.0.0.1:7474/")
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

func TestQuerySimple(t *testing.T) {
	rows := prepareAndQuery("return 1")
	hasNext := rows.Next()
	if !hasNext {
		t.Fatal("no next!")
	}

	var test int
	err := rows.Scan(&test)
	if err != nil {
		t.Fatal(err)
	}

	if test != 1 {
		t.Fatal("test != 1")
	}
}

func TestQuerySimpleFloat(t *testing.T) {
	rows := prepareAndQuery("return 1.2")
	rows.Next()
	var test float64
	err := rows.Scan(&test)
	if err != nil {
		t.Fatal(err)
	}

	if test != 1.2 {
		t.Fatal("test != 1.2")
	}
}

func TestQuerySimpleString(t *testing.T) {
	rows := prepareAndQuery("return '123'")
	rows.Next()
	var test string
	err := rows.Scan(&test)
	if err != nil {
		t.Fatal(err)
	}

	if test != "123" {
		t.Fatal("test != '123'")
	}
}

func TestQueryIntParam(t *testing.T) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(123)
	if err != nil {
		t.Fatal(err)
	}
	rows.Next()
	var test int
	err = rows.Scan(&test)
	if err != nil {
		t.Fatal(err)
	}
	if test != 123 {
		t.Fatal("test != 123;", test)
	}
}

/*
func TestTransactionRollback(t *testing.T) {
	db := testConn()
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := tx.Query("return 1")
	if err != nil {
		t.Fatal(err)
	}

	hasNext := rows.Next()
	if !hasNext {
		t.Fatal("no rows in transaction")
	}

	var test int
	rows.Scan(&test)
	if test != 1 {
		t.Fatal("test != 1")
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(err)
	}
}*/

// TODO array conversion
/*
func TestQuerySimpleIntArray(t *testing.T) {
	rows := prepareAndQuery("return [1,2,3]")
	rows.Next()
	var test []int
	err := rows.Scan(&test)
	if err != nil {
		t.Fatal(err)
	}

	if test[0] != 1 || test[1] != 2 || test[2] != 3 {
		t.Fatal("test != [1,2,3];", test)
	}
} */
