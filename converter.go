package cq

import (
	"database/sql/driver"
	//	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
)

var _ driver.ValueConverter = cypherType{}

type cypherType struct {
}

func (i cypherType) Scan(value interface{}) error {
	log.Println("Scan", value)
	switch d := value.(type) {
	case *[]int:
		*d = []int{1, 2, 3}
	}
	return errors.New("Scan failed")
}

func (i cypherType) Value() (driver.Value, error) {
	log.Println("Value", i)
	return []byte("[1,2,3]"), nil
}

func (i cypherType) ConvertValue(v interface{}) (driver.Value, error) {
	if svi, ok := v.(driver.Valuer); ok {
		_, err := svi.Value()
		if err != nil {
			return nil, err
		}
	}
	//log.Println("Converting Value", v)
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int:
		return rv.Int(), nil
	case reflect.String:
		return rv.String(), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	}
	return nil, errors.New(fmt.Sprintf("cq: unsupported value %v (type %T) converting to []int", v, v))
}

func (cs cypherStmt) ColumnConverter(idx int) driver.ValueConverter {
	return cypherType{}
}
