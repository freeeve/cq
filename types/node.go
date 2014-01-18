package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Node struct {
	LabelURI   string                 `json:"labels"`
	SelfURI    string                 `json:"self"`
	Properties map[string]CypherValue `json:"data"`
}

func (n *Node) Scan(value interface{}) error {
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case []byte:
		err := json.Unmarshal(value.([]byte), &n)
		return err
	}
	return errors.New("cq: invalid Scan value for Node")
}

func (n *Node) Labels() ([]string, error) {
	fmt.Println(n)
	resp, err := http.Get(n.LabelURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []string{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	return ret, err
}
