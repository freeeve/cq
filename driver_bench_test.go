package cq_test

import (
	"log"
	"testing"
)

func BenchmarkSimpleQuery(b *testing.B) {
	stmt := prepareTest("return 1")
	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		var test int
		rows.Scan(&test)
	}
}

func BenchmarkSimpleCreate(b *testing.B) {
	stmt := prepareTest("create ()")
	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		var test int
		rows.Scan(&test)
	}
}

func BenchmarkSimpleCreateLabel(b *testing.B) {
	stmt := prepareTest("create (:Test)")
	for i := 0; i < b.N; i++ {
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		var test int
		rows.Scan(&test)
	}
}
