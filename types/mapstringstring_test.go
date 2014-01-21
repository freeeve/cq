package types_test

import (
	"errors"
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

func (s *TypesSuite) TestQueryMapStringStringParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(
		types.MapStringString{map[string]string{"key1": "1", "key2": "2"}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.MapStringString
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, map[string]string{"key1": "1", "key2": "2"})
}

func (s *TypesSuite) TestQueryStringStringMapParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(map[string]string{"key1": "1", "key2": "2"})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.MapStringString
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, map[string]string{"key1": "1", "key2": "2"})
}

func (s *TypesSuite) TestQueryMapStringString(c *C) {
	rows := prepareAndQuery(`return {key1:"1",key2:"2"}`)
	rows.Next()
	var test types.MapStringString
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val, DeepEquals, map[string]string{"key1": "1", "key2": "2"})
}

func (s *TypesSuite) TestQueryBadMapStringString(c *C) {
	rows := prepareAndQuery(`return {key1:1}`)
	rows.Next()
	var test types.MapStringString
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal number into Go value of type string"))
}

func (s *TypesSuite) TestQueryNullMapStringString(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.MapStringString
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}
