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

func BenchmarkTransactional10SimpleCreate(b *testing.B) {
	transactionalSizeSimpleCreate(b, 10)
}

func BenchmarkTransactional100SimpleCreate(b *testing.B) {
	transactionalSizeSimpleCreate(b, 100)
}

func BenchmarkTransactional1000SimpleCreate(b *testing.B) {
	transactionalSizeSimpleCreate(b, 1000)
}

func BenchmarkTransactional10000SimpleCreate(b *testing.B) {
	transactionalSizeSimpleCreate(b, 10000)
}

func transactionalSizeSimpleCreate(b *testing.B, size int) {
	conn := testConn()
	defer conn.Close()

	tx, err := conn.Begin()
	//	log.Println("begin:", tx)
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("create ({n:{0}})")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		//		log.Println("i:", i)
		_, err = stmt.Exec(i)
		if err != nil {
			log.Fatal(err)
		}
		if (i > 0 && i%size == 0) || i == b.N-1 {
			//			log.Println("committing:", tx)
			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			}
			if i < b.N-1 {
				tx, err = conn.Begin()
				//				log.Println("begin:", tx)
				if err != nil {
					log.Fatal(err)
				}
				stmt, err = tx.Prepare("create ({n:{0}})")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
