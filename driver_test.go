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

func TestQuerySimple(t *testing.T) {
	db := testConn()
	stmt, err := db.Prepare("return 1")
	if err != nil {
		t.Fatal(err)
	}

	rows, err := stmt.Query()
	if err != nil {
		t.Fatal(err)
	}

	hasNext := rows.Next()
	if !hasNext {
		t.Fatal("no next!")
	}

	var test int
	err = rows.Scan(&test)
	if err != nil {
		t.Fatal(err)
	}

	if test != 1 {
		t.Fatal("test != 1")
	}

}
