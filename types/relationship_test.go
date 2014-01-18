package types_test

import (
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

func (s *TypesSuite) TestQueryRelationship(c *C) {
	stmt := prepareTest(`create (:Test)-[r:TEST_TYPE {foo:"bar", i:1}]->(:Test) return r`)
	rows, err := stmt.Query()
	c.Assert(err, IsNil)

	rows.Next()
	var test types.Relationship
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	t1 := types.Relationship{}
	t1.Properties = map[string]types.CypherValue{}
	t1.Properties["foo"] = types.CypherValue{"bar", types.CypherString}
	t1.Properties["i"] = types.CypherValue{1, types.CypherInt}
	c.Assert(test.Properties, DeepEquals, t1.Properties)
	c.Assert(test.Type, Equals, "TEST_TYPE")
}
