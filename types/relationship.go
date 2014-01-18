package types

import (
	"encoding/json"
	"errors"
)

type Relationship struct {
	Type       string                 `json:"type"`
	SelfURI    string                 `json:"self"`
	Properties map[string]CypherValue `json:"data"`
}

func (r *Relationship) Scan(value interface{}) error {
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case []byte:
		err := json.Unmarshal(value.([]byte), &r)
		return err
	}
	return errors.New("cq: invalid Scan value for Relationship")
}
