package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ArrayString struct {
	Val []string
}

func (as *ArrayString) Scan(value interface{}) error {
	fmt.Println("attempting to Scan:", value)
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case string:
		err := json.Unmarshal([]byte(value.(string)), &as.Val)
		return err
	case []byte:
		err := json.Unmarshal(value.([]byte), &as.Val)
		return err
	}
	return errors.New("cq: invalid Scan value for ArrayString")
}

func (as ArrayString) Value() (driver.Value, error) {
	b, err := json.Marshal(as.Val)
	fmt.Println("valued:", string(b))
	return string(b), err
}
