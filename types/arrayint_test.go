package types_test

import (
	"errors"
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

func (s *TypesSuite) TestQueryArrayIntParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayInt{[]int{1, 2, 3}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int{1, 2, 3})
}

func (s *TypesSuite) TestQueryIntArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]int{1, 2, 3})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int{1, 2, 3})
}

func (s *TypesSuite) TestQueryArrayInt(c *C) {
	rows := prepareAndQuery("return [1,2,3]")
	rows.Next()
	var test types.ArrayInt
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int{1, 2, 3})
}

func (s *TypesSuite) TestQueryBadIntArray(c *C) {
	rows := prepareAndQuery("return [1,2,'asdf']")
	rows.Next()
	var test types.ArrayInt
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type int"))
}

func (s *TypesSuite) TestQueryNullIntArray(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayInt
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}

func (s *TypesSuite) TestQueryArrayInt64Param(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayInt64{[]int64{12345678910, 234567891011, 3456789101112}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int64{12345678910, 234567891011, 3456789101112})
}

func (s *TypesSuite) TestQueryInt64ArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]int64{12345678910, 234567891011, 3456789101112})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInt64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int64{12345678910, 234567891011, 3456789101112})
}

func (s *TypesSuite) TestQueryArrayInt64(c *C) {
	rows := prepareAndQuery("return [12345678910, 234567891011, 3456789101112]")
	rows.Next()
	var test types.ArrayInt64
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []int64{12345678910, 234567891011, 3456789101112})
}

func (s *TypesSuite) TestQueryBadInt64Array(c *C) {
	rows := prepareAndQuery("return [123456789,'asdf']")
	rows.Next()
	var test types.ArrayInt64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type int64"))
}

func (s *TypesSuite) TestQueryNullInt64Array(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayInt64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}
