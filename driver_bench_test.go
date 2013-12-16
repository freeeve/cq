package cq_test

import (
	. "launchpad.net/gocheck"
	"log"
)

type BenchmarkSuite struct{}

var _ = Suite(&BenchmarkSuite{})

func (s *BenchmarkSuite) SetUpTest(c *C) {
	db := testConn()
	db.Exec("match (n) where has(n.`benchmark~test~id`) delete n")
}

func (s *BenchmarkSuite) TearDownTest(c *C) {
	db := testConn()
	db.Exec("match (n) where has(n.`benchmark~test~id`) delete n")
}

func (s *BenchmarkSuite) BenchmarkSimpleQuery(c *C) {
	stmt := prepareTest("return 1")
	defer stmt.Close()
	var test int
	for i := 0; i < c.N; i++ {
		rows, err := stmt.Query()
		c.Assert(err, IsNil)

		rows.Scan(&test)
		rows.Close()
	}
}

func (s *BenchmarkSuite) BenchmarkSimpleCreate(c *C) {
	stmt := prepareTest("create ({`benchmark~test~id`:0})")
	defer stmt.Close()
	var test int
	for i := 0; i < c.N; i++ {
		rows, err := stmt.Query()
		c.Assert(err, IsNil)

		rows.Scan(&test)
		rows.Close()
	}
}

func (s *BenchmarkSuite) BenchmarkSimpleCreateLabel(c *C) {
	stmt := prepareTest("create (:Test {`benchmark~test~id`:0})")
	defer stmt.Close()
	var test int
	for i := 0; i < c.N; i++ {
		rows, err := stmt.Query()
		c.Assert(err, IsNil)

		rows.Scan(&test)
		rows.Close()
	}
}

func (s *BenchmarkSuite) BenchmarkTx10SimpleCreate(c *C) {
	txSizeSimpleCreate(c, 10)
}

func (s *BenchmarkSuite) BenchmarkTx100SimpleCreate(c *C) {
	txSizeSimpleCreate(c, 100)
}

func (s *BenchmarkSuite) BenchmarkTx1000SimpleCreate(c *C) {
	txSizeSimpleCreate(c, 1000)
}

func (s *BenchmarkSuite) BenchmarkTx10000SimpleCreate(c *C) {
	txSizeSimpleCreate(c, 10000)
}

func txSizeSimpleCreate(c *C, size int) {
	conn := testConn()
	defer conn.Close()

	tx, err := conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("create ({`benchmark~test~id`:{0}})")
	c.Assert(err, IsNil)

	for i := 0; i < c.N; i++ {
		_, err = stmt.Exec(i)
		c.Assert(err, IsNil)

		if (i > 0 && i%size == 0) || i == c.N-1 {
			err = tx.Commit()
			c.Assert(err, IsNil)

			if i < c.N-1 {
				tx, err = conn.Begin()
				c.Assert(err, IsNil)

				stmt.Close()
				stmt, err = tx.Prepare("create ({`benchmark~test~id`:{0}})")
				c.Assert(err, IsNil)
			}
		}
	}
}
