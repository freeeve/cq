package cq_test

import (
	"database/sql"
	"errors"
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
	"log"
)

// This file is meant to hold integration tests where cq must be imported

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

func (s *DriverSuite) TestQueryFloatParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(1234567910.891)
	rows.Next()
	var test float64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test, Equals, 1234567910.891)
}

func (s *DriverSuite) TestQuerySimpleString(c *C) {
	rows := prepareAndQuery("return '123'")
	rows.Next()
	var test string
	err := rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != "123" {
		c.Fatal("test != '123';", test)
	}
}

func (s *DriverSuite) TestQueryStringParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query("123")
	rows.Next()
	var test string
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test, Equals, "123")
}

func (s *DriverSuite) TestQueryArrayByteParam(c *C) {
	c.Skip("byte arrays don't work yet")
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]byte("123"))
	rows.Next()
	var test []byte
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(string(test), DeepEquals, string([]byte("123")))
}

func (s *DriverSuite) TestQuerySimpleBool(c *C) {
	rows := prepareAndQuery("return true")
	rows.Next()
	var test bool
	err := rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != true {
		c.Fatal("test != true;", test)
	}
}

func (s *DriverSuite) TestQueryBoolParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(true)
	rows.Next()
	var test bool
	err = rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != true {
		c.Fatal("test != true;", test)
	}
}

func (s *DriverSuite) TestQueryBoolFalseParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(false)
	rows.Next()
	var test bool
	err = rows.Scan(&test)
	c.Assert(err, IsNil)

	if test != false {
		c.Fatal("test != true;", test)
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
	c.Assert(test, Equals, 123)
}

func (s *DriverSuite) TestQueryArrayIntParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayInt{[]int{1, 2, 3}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int{1, 2, 3})
}

func (s *DriverSuite) TestQueryIntArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]int{1, 2, 3})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int{1, 2, 3})
}

func (s *DriverSuite) TestQueryArrayInt(c *C) {
	rows := prepareAndQuery("return [1,2,3]")
	rows.Next()
	var test types.ArrayInt
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int{1, 2, 3})
}

func (s *DriverSuite) TestQueryBadIntArray(c *C) {
	rows := prepareAndQuery("return [1,2,'asdf']")
	rows.Next()
	var test types.ArrayInt
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type int"))
}

func (s *DriverSuite) TestQueryNullIntArray(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayInt
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}

func (s *DriverSuite) TestQueryArrayFloat64Param(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayFloat64{[]float64{1.1, 2.1, 3.1}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayFloat64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []float64{1.1, 2.1, 3.1})
}

func (s *DriverSuite) TestQueryFloat64ArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]float64{1.1, 2.1, 3.1})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayFloat64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []float64{1.1, 2.1, 3.1})
}

func (s *DriverSuite) TestQueryArrayFloat64(c *C) {
	rows := prepareAndQuery("return [1.1,2.1,3.1]")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []float64{1.1, 2.1, 3.1})
}

func (s *DriverSuite) TestQueryBadFloatArray(c *C) {
	rows := prepareAndQuery("return [1.1,2.1,'asdf']")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type float64"))
}

func (s *DriverSuite) TestQueryNullFloat64Array(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}

func (s *DriverSuite) TestQueryArrayInt64Param(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayInt64{[]int64{12345678910, 234567891011, 3456789101112}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int64{12345678910, 234567891011, 3456789101112})
}

func (s *DriverSuite) TestQueryInt64ArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]int64{12345678910, 234567891011, 3456789101112})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int64{12345678910, 234567891011, 3456789101112})
}

func (s *DriverSuite) TestQueryArrayInt64(c *C) {
	rows := prepareAndQuery("return [12345678910, 234567891011, 3456789101112]")
	rows.Next()
	var test types.ArrayInt64
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int64{12345678910, 234567891011, 3456789101112})
}

func (s *DriverSuite) TestQueryBadInt64Array(c *C) {
	rows := prepareAndQuery("return [123456789,'asdf']")
	rows.Next()
	var test types.ArrayInt64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type int64"))
}

func (s *DriverSuite) TestQueryNullInt64Array(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayInt64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}

func (s *DriverSuite) TestQueryCypherValueNull(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.CypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Type, Equals, types.CypherNull)
	c.Assert(test.Val, Equals, nil)
}

func (s *DriverSuite) TestQueryCypherValueBoolean(c *C) {
	rows := prepareAndQuery("return true")
	rows.Next()
	var test types.CypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Type, Equals, types.CypherBoolean)
	c.Assert(test.Val, Equals, true)
}

func (s *DriverSuite) TestQueryCypherValueString(c *C) {
	rows := prepareAndQuery("return 'asdf'")
	rows.Next()
	var test types.CypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Type, Equals, types.CypherString)
	c.Assert(test.Val, Equals, "asdf")
}

func (s *DriverSuite) TestQueryCypherValueInt64(c *C) {
	rows := prepareAndQuery("return 9223372000000000000")
	rows.Next()
	var test types.CypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, Equals, int64(9223372000000000000))
	c.Assert(test.Type, Equals, types.CypherInt64)
}

func (s *DriverSuite) TestQueryCypherValueInt(c *C) {
	rows := prepareAndQuery("return 1234567890")
	rows.Next()
	var test types.CypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Type, Equals, types.CypherInt)
	c.Assert(test.Val, Equals, 1234567890)
}

func (s *DriverSuite) TestQueryCypherValueIntArray(c *C) {
	rows := prepareAndQuery("return [1,2,2345678910]")
	rows.Next()
	var test types.CypherValue
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Type, Equals, types.CypherArrayInt)
	c.Assert(test.Val.([]int), DeepEquals, []int{1, 2, 2345678910})
}

func (s *DriverSuite) TestQueryNullString(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var nullString sql.NullString
	err := rows.Scan(&nullString)
	c.Assert(err, IsNil)
	c.Assert(nullString.Valid, Equals, false)
}

func (s *DriverSuite) TestScanNullInt64(c *C) {
	rows := prepareAndQuery("return 123456789")
	rows.Next()
	var nullInt64 sql.NullInt64
	err := rows.Scan(&nullInt64)
	c.Assert(err, IsNil)
	c.Assert(nullInt64.Valid, Equals, true)
	c.Assert(nullInt64.Int64, Equals, int64(123456789))
}

func (s *DriverSuite) TestScanBigInt64(c *C) {
	rows := prepareAndQuery("return 123456789101112")
	rows.Next()
	var i64 int64
	err := rows.Scan(&i64)
	c.Assert(err, IsNil)
	c.Assert(i64, Equals, int64(123456789101112))
}

func (s *DriverSuite) TestExecNilRows(c *C) {
	db := testConn()
	db.Exec("...")
}
