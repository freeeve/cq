package types

import (
	"fmt"
	"time"
)

type CypherTime struct {
	Time  time.Time
	Valid bool
}

func (ct *CypherTime) Scan(value interface{}) error {
	if value == nil {
		ct.Valid = false
		return nil
	}

	switch value.(type) {
	// do we need to handle int64 too?
	case int:
		ct.Time = time.Unix(0, int64(value.(int)*1000000))
		ct.Valid = true
		return nil
	case CypherValue:
		cv := value.(CypherValue)
		if cv.Type == CypherInt64 {
			ct.Time = time.Unix(0, cv.Val.(int64)*1000000)
			ct.Valid = true
			return nil
		}
	default:
		fmt.Println(value)
	}
	ct.Valid = false
	return nil
}
