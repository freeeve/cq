package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ArrayInt struct {
	Val []int
}

func (ai *ArrayInt) Scan(value interface{}) error {
	//fmt.Println("attempting to Scan:", value)
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case string:
		err := json.Unmarshal([]byte(value.(string)), &ai.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &ai.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayInt")
}

func (ai ArrayInt) Value() (driver.Value, error) {
	b, err := json.Marshal(ai.Val)
	return string(b), err
}

type ArrayInt64 struct {
	Val []int64
}

func (ai *ArrayInt64) Scan(value interface{}) error {
	//fmt.Println("attempting to Scan:", value)
	if value == nil {
		return errors.New("cq: scan value is null")
	}

	switch value.(type) {
	case string:
		err := json.Unmarshal([]byte(value.(string)), &ai.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &ai.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayInt")
}

func (ai ArrayInt64) Value() (driver.Value, error) {
	b, err := json.Marshal(ai.Val)
	return string(b), err
}
