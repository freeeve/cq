package cq

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CypherType uint8

// supported types
const (
	CypherNull               CypherType = iota
	CypherBoolean            CypherType = iota
	CypherString             CypherType = iota
	CypherInt64              CypherType = iota
	CypherInt                CypherType = iota
	CypherFloat64            CypherType = iota
	CypherArrayInt           CypherType = iota
	CypherArrayInt64         CypherType = iota
	CypherArrayByte          CypherType = iota
	CypherArrayFloat64       CypherType = iota
	CypherArrayString        CypherType = iota
	CypherArrayInterface     CypherType = iota
	CypherMapStringString    CypherType = iota
	CypherMapStringInterface CypherType = iota
)

type ArrayInt struct {
	Value []int
}

func (ai *ArrayInt) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	err := json.Unmarshal(value.([]byte), &ai.Value)
	return err
}

func (v *CypherValue) Scan(value interface{}) error {
	if v == nil {
		return ErrScanOnNil
	}
	if value == nil {
		v.Value = nil
		v.Type = CypherNull
		return nil
	}

	switch value.(type) {
	case bool:
		v.Type = CypherBoolean
		v.Value = value
		return nil
	case string:
		v.Type = CypherString
		v.Value = value
		return nil
	case int:
		if value.(int) > ((1 << 31) - 1) {
			v.Type = CypherInt64
			v.Value = int64(value.(int))
			return nil
		}
		v.Type = CypherInt
		v.Value = value
		return nil
	}

	err := json.Unmarshal(value.([]byte), &v)
	if err != nil {
		return err
	}

	switch v.Type {
	case CypherArrayInt:
		var ai ArrayInt
		err = json.Unmarshal(value.([]byte), &ai.Value)
		v.Value = ai.Value
		return err
	}
	return err
}

type CypherValue struct {
	Value interface{}
	Type  CypherType
}

func (c *CypherValue) UnmarshalJSON(b []byte) error {
	str := string(b)
	switch str {
	case "null":
		c.Value = nil
		c.Type = CypherNull
		return nil
	case "true":
		c.Value = true
		c.Type = CypherBoolean
		return nil
	case "false":
		c.Value = false
		c.Type = CypherBoolean
		return nil
	}
	if len(b) > 0 {
		if b[0] == byte('"') {
			c.Value = strings.Trim(str, "\"")
			c.Type = CypherString
			return nil
		}
	}
	var err error
	c.Value, err = strconv.Atoi(str)
	if err == nil {
		c.Type = CypherInt
		return nil
	}
	c.Value, err = strconv.ParseInt(str, 10, 64)
	if err == nil {
		c.Type = CypherInt64
		return nil
	}
	c.Value, err = strconv.ParseFloat(str, 64)
	if err == nil {
		c.Type = CypherFloat64
		return nil
	}
	c.Value = b
	c.Type = CypherArrayInt
	//json.Unmarshal(b, &c.Value)
	return nil
}

func (CypherValue) ConvertValue(v interface{}) (driver.Value, error) {
	if driver.IsValue(v) {
		return v, nil
	}

	if svi, ok := v.(driver.Valuer); ok {
		sv, err := svi.Value()
		if err != nil {
			return nil, err
		}
		if !driver.IsValue(sv) {
			return nil, fmt.Errorf("non-Value type %T returned from Value", sv)
		}
		return sv, nil
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr:
		// indirect pointers
		if rv.IsNil() {
			return nil, nil
		} else {
			return CypherValue{}.ConvertValue(rv.Elem().Interface())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(rv.Uint()), nil
	case reflect.Uint64:
		u64 := rv.Uint()
		if u64 >= 1<<63 {
			return nil, fmt.Errorf("uint64 values with high bit set are not supported")
		}
		return int64(u64), nil
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	}
	return nil, fmt.Errorf("unsupported type %T, a %s", v, rv.Kind())
}

func (cs cypherStmt) ColumnConverter(idx int) driver.ValueConverter {
	return CypherValue{}
}
