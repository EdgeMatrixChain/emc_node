package jsonrpc

import (
	"github.com/emc-protocol/edge-matrix/application"
)

// NodeQuery is a query to filter node
type NodeQuery struct {
	name    string
	tag     string
	id      string
	version string
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
