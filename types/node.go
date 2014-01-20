package types

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Conn interface {
	BaseURL() string
}

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

func (n *Node) Labels(baseURL string) ([]string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	labelURL, err := url.Parse(n.LabelURI)
	if err != nil {
		return nil, err
	}
	labelURL.Scheme = base.Scheme
	labelURL.User = base.User
	resp, err := http.Get(n.LabelURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []string{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	return ret, err
}
