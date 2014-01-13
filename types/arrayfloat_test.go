package types_test

import (
	"errors"
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

func (s *TypesSuite) TestQueryArrayFloat64Param(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayFloat64{[]float64{1.1, 2.1, 3.1}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayFloat64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []float64{1.1, 2.1, 3.1})
}

func (s *TypesSuite) TestQueryFloat64ArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query([]float64{1.1, 2.1, 3.1})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayFloat64
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []float64{1.1, 2.1, 3.1})
}

func (s *TypesSuite) TestQueryArrayFloat64(c *C) {
	rows := prepareAndQuery("return [1.1,2.1,3.1]")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, []float64{1.1, 2.1, 3.1})
}

func (s *TypesSuite) TestQueryBadFloatArray(c *C) {
	rows := prepareAndQuery("return [1.1,2.1,'asdf']")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type float64"))
}

func (s *TypesSuite) TestQueryNullFloat64Array(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}
