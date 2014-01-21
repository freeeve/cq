package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
)

type MapStringString struct {
	Val map[string]string
}

func (mss *MapStringString) Scan(value interface{}) error {
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
		err = json.Unmarshal([]byte(str), &mss.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &mss.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayInterface")
}

func (mss MapStringString) Value() (driver.Value, error) {
	b, err := json.Marshal(mss.Val)
	return string(b), err
}
