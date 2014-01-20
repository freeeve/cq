package types_test

import (
	"errors"
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

type MyStruct struct {
	Foo string
	Bar int
}

func (s *TypesSuite) TestQueryArrayInterfaceParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(types.ArrayInterface{
		[]interface{}{
			MyStruct{"1", 1},
			MyStruct{"2", 2},
			MyStruct{"3", 3},
		}})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInterface
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val[0], DeepEquals,
		map[string]interface{}{"Foo": "1", "Bar": 1.0})
	c.Assert(test.Val[1], DeepEquals,
		map[string]interface{}{"Foo": "2", "Bar": 2.0})
	c.Assert(test.Val[2], DeepEquals,
		map[string]interface{}{"Foo": "3", "Bar": 3.0})
}

func (s *TypesSuite) TestQueryInterfaceArrayParam(c *C) {
	stmt := prepareTest("with {0} as test return test")
	rows, err := stmt.Query(
		[]MyStruct{
			MyStruct{"1", 1},
			MyStruct{"2", 2},
			MyStruct{"3", 3},
		})
	c.Assert(err, IsNil)

	rows.Next()
	var test types.ArrayInterface
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val[0], DeepEquals,
		map[string]interface{}{"Foo": "1", "Bar": 1.0})
	c.Assert(test.Val[1], DeepEquals,
		map[string]interface{}{"Foo": "2", "Bar": 2.0})
	c.Assert(test.Val[2], DeepEquals,
		map[string]interface{}{"Foo": "3", "Bar": 3.0})
}

func (s *TypesSuite) TestQueryArrayInterface(c *C) {
	rows := prepareAndQuery(`return [{Foo:"1",Bar:1},{Foo:"2",Bar:2},{Foo:"3",Bar:3}]`)
	rows.Next()
	var test types.ArrayInterface
	err := rows.Scan(&test)
	c.Assert(err, IsNil)
	c.Assert(test.Val[0], DeepEquals,
		map[string]interface{}{"Foo": "1", "Bar": 1.0})
	c.Assert(test.Val[1], DeepEquals,
		map[string]interface{}{"Foo": "2", "Bar": 2.0})
	c.Assert(test.Val[2], DeepEquals,
		map[string]interface{}{"Foo": "3", "Bar": 3.0})
}

func (s *TypesSuite) TestQueryBadInterfaceArray(c *C) {
	rows := prepareAndQuery("return [1.1,2.1,'asdf']")
	rows.Next()
	var test types.ArrayFloat64
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: json: cannot unmarshal string into Go value of type float64"))
}

func (s *TypesSuite) TestQueryNullInterfaceArray(c *C) {
	rows := prepareAndQuery("return null")
	rows.Next()
	var test types.ArrayInterface
	err := rows.Scan(&test)
	c.Assert(err, DeepEquals, errors.New("sql: Scan error on column index 0: cq: scan value is null"))
}
