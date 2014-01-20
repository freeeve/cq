package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type ArrayString struct {
	Val []string
}

func (as *ArrayString) Scan(value interface{}) error {
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
		err = json.Unmarshal([]byte(str), &as.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &as.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayString")
}

func (as ArrayString) Value() (driver.Value, error) {
	b, err := json.Marshal(as.Val)
	return string(b), err
}
