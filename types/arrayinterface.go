package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type ArrayInterface struct {
	Val []interface{}
}

func (ai *ArrayInterface) Scan(value interface{}) error {
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
		err = json.Unmarshal([]byte(str), &ai.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &ai.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayInterface")
}

func (ai ArrayInterface) Value() (driver.Value, error) {
	b, err := json.Marshal(ai.Val)
	return string(b), err
}
