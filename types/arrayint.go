package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type ArrayInt struct {
	Val []int
}

func (ai *ArrayInt) Scan(value interface{}) error {
	if value == nil {
		return ErrScanOnNil
	}

	/*
		cv := CypherValue{}
		err := json.Unmarshal(value.([]byte), &cv)
		if err != nil {
			return err
		}
		if cv.Type == CypherArrayInt {
			ai.Val = cv.Val.([]int)
			return nil
		}
	*/
	err := json.Unmarshal(value.([]byte), &ai.Val)
	if err != nil {
		return err
	}
	return nil
	return errors.New("cq: invalid Scan value for ArrayInt")
}

func (ai ArrayInt) Value() (driver.Value, error) {
	b, err := json.Marshal(CypherValue{CypherArrayInt, ai.Val})
	fmt.Println("Value(): ", string(b))
	return b, err
}

type ArrayInt64 struct {
	Val []int64
}

func (ai *ArrayInt64) Scan(value interface{}) error {
	//fmt.Println("attempting to Scan:", value)
	if value == nil {
		return errors.New("cq: scan value is null")
	}

	/*
		cv := CypherValue{}
		err := json.Unmarshal(value.([]byte), &cv)
		if err != nil {
			return err
		}
		ai.Val = cv.Val.([]int64)
	*/
	err := json.Unmarshal(value.([]byte), &ai.Val)
	if err != nil {
		return err
	}
	return nil
	return errors.New("cq: invalid Scan value for ArrayInt")
}

func (ai ArrayInt64) Value() (driver.Value, error) {
	b, err := json.Marshal(CypherValue{CypherArrayInt, ai.Val})
	fmt.Println("Value(): ", string(b))
	return b, err
}
