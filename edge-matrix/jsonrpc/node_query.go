package jsonrpc

import (
	"encoding/json"
	"github.com/emc-protocol/edge-matrix/application"
)

// NodeQuery is a query to filter node
type NodeQuery struct {
	name    string
	tag     string
	id      string
	version string
}

func decodeNodeQueryFromInterface(i interface{}) (*NodeQuery, error) {
	// once the node filter is decoded as map[string]interface we cannot use unmarshal json
	raw, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	query := &NodeQuery{}
	if err := json.Unmarshal(raw, &query); err != nil {
		return nil, err
	}

	return query, nil
}

func (q *NodeQuery) Match(rm *application.Application) bool {
	if q.tag != "" {
		match := false
		if q.tag == rm.Tag {
			match = true
		}
		if !match {
			return false
		}
	}
	// check name
	if q.name != "" {
		match := false
		if rm.Name == q.name {
			match = true
		}

		if !match {
			return false
		}
	}

	if q.id != "" {
		match := false
		if rm.PeerID.String() == q.id {
			match = true
		}

		if !match {
			return false
		}
	}

	if q.version != "" {
		match := false
		if rm.Version == q.version {
			match = true
		}

		if !match {
			return false
		}
	}

	return true
}
