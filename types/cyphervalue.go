package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CypherType uint8

var (
	ErrScanOnNil = errors.New("cq: scan value is null")
)

// supported types
const (
	CypherNull            CypherType = iota
	CypherBoolean         CypherType = iota
	CypherString          CypherType = iota
	CypherInt64           CypherType = iota
	CypherInt             CypherType = iota
	CypherFloat64         CypherType = iota
	CypherArrayInt        CypherType = iota
	CypherArrayInt64      CypherType = iota
	CypherArrayByte       CypherType = iota
	CypherArrayFloat64    CypherType = iota
	CypherArrayString     CypherType = iota
	CypherArrayInterface  CypherType = iota
	CypherMapStringString CypherType = iota
	CypherNode            CypherType = iota
	CypherRelationship    CypherType = iota
	CypherPath            CypherType = iota
)

func (v *CypherValue) Scan(value interface{}) error {
	//fmt.Println("attempting to Scan:", value)
	if v == nil {
		return ErrScanOnNil
	}
	if value == nil {
		v.Val = nil
		v.Type = CypherNull
		return nil
	}

	switch value.(type) {
	case bool:
		v.Type = CypherBoolean
		v.Val = value
		return nil
	case string:
		v.Type = CypherString
		v.Val = value
		return nil
	case int:
		if value.(int) > ((1 << 31) - 1) {
			v.Type = CypherInt64
			v.Val = int64(value.(int))
			return nil
		}
		v.Type = CypherInt
		v.Val = value
		return nil
	}

	err := json.Unmarshal(value.([]byte), &v)
	if err != nil {
		return err
	}

	switch v.Type {
	case CypherArrayInt:
		var ai ArrayInt
		err = json.Unmarshal(value.([]byte), &ai.Val)
		v.Val = ai.Val
		return err
	}
	return err
}

type CypherValue struct {
	Val  interface{}
	Type CypherType
}

/*
func (cv *CypherValue) Value() (driver.Value, error) {
	fmt.Println(cv, "Value()")
	fmt.Println(cv.Val)
	b, err := json.Marshal(cv)
	return string(b), err
}*/

func (c *CypherValue) UnmarshalJSON(b []byte) error {
	//fmt.Println("attempting to unmarshal: ", string(b))
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err == nil {
		if m["Type"] != nil {
			c.Val = m["Val"]
			c.Type = m["Type"].(CypherType)
			return nil
		}
	}
	err = nil
	str := string(b)
	switch str {
	case "null":
		c.Val = nil
		c.Type = CypherNull
		return nil
	case "true":
		c.Val = true
		c.Type = CypherBoolean
		return nil
	case "false":
		c.Val = false
		c.Type = CypherBoolean
		return nil
	}
	if len(b) > 0 {
		if b[0] == byte('"') {
			c.Val = strings.Trim(str, "\"")
			c.Type = CypherString
			return nil
		}
	}
	c.Val, err = strconv.Atoi(str)
	if err == nil {
		c.Type = CypherInt
		return nil
	}
	c.Val, err = strconv.ParseInt(str, 10, 64)
	if err == nil {
		c.Type = CypherInt64
		return nil
	}
	c.Val, err = strconv.ParseFloat(str, 64)
	if err == nil {
		c.Type = CypherFloat64
		return nil
	}
	c.Val = b
	c.Type = CypherArrayInt
	//json.Unmarshal(b, &c.Val)
	return nil
}

func (cv CypherValue) ConvertValue(v interface{}) (driver.Value, error) {
	//fmt.Println("attempting to convert:", v)
	if driver.IsValue(v) {
		//fmt.Println("IsValue")
		return v, nil
	}

	if svi, ok := v.(driver.Valuer); ok {
		//fmt.Println("we have a valuer:", v)
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
	case reflect.Slice:
		b, err := json.Marshal(v)
		return string(b), err
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
