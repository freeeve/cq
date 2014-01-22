package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type MapStringCypherValue struct {
	Val map[string]CypherValue
}

func (msc *MapStringCypherValue) Scan(value interface{}) error {
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case string:
		str := "\"" + value.(string) + "\""
		str, err := strconv.Unquote(str)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(str), &msc.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &msc.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for MapStringCypherValue")
}

func (msc MapStringCypherValue) Value() (driver.Value, error) {
	b, err := json.Marshal(msc.Val)
	return string(b), err
}
