package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ArrayFloat64 struct {
	Val []float64
}

func (af *ArrayFloat64) Scan(value interface{}) error {
	//fmt.Println("attempting to Scan:", value)
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case string:
		err := json.Unmarshal([]byte(value.(string)), &af.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &af.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayFloat")
}

func (ai ArrayFloat64) Value() (driver.Value, error) {
	b, err := json.Marshal(ai.Val)
	return string(b), err
}
