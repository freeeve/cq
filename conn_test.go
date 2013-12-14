package cq

import (
	"database/sql/driver"
	"log"
	"testing"
)

var (
	testURL = "http://localhost:7474/"
)

func openTest() driver.Conn {
	db, err := Open(testURL)
	if err != nil {
		log.Println("can't connect to db.")
		return nil
	}
	return db
}

func TestOpen(t *testing.T) {
	db := openTest()
	if db == nil {
		t.Fatal("can't connect to test db: ", testURL)
	}
}

func TestPrepareNoParams(t *testing.T) {
	db := openTest()
	if db == nil {
		t.Fatal("can't connect to test db: ", testURL)
	}
	stmt, err := db.Prepare("match (n) return n limit 1")
	if err != nil {
		t.Fatal(err)
	}
	if stmt == nil {
		t.Fatal("statement shouldn't be nil")
	}
}
