package types_test

import (
	_ "github.com/wfreeman/cq"
	"github.com/wfreeman/cq/types"
	. "launchpad.net/gocheck"
)

func (s *TypesSuite) TestQueryNode(c *C) {
	stmt := prepareTest(`create (a:Test {foo:"bar", i:1}) return a`)
	rows, err := stmt.Query()
	c.Assert(err, IsNil)

	rows.Next()
	var test types.Node
	err = rows.Scan(&test)
	c.Assert(err, IsNil)
	t1 := types.Node{}
	t1.Properties = map[string]types.CypherValue{}
	t1.Properties["foo"] = types.CypherValue{"bar", types.CypherString}
	t1.Properties["i"] = types.CypherValue{1, types.CypherInt}
	c.Assert(test.Properties, DeepEquals, t1.Properties)
	labels, err := test.Labels()
	c.Assert(err, IsNil)
	c.Assert(labels, DeepEquals, []string{"Test"})
}
