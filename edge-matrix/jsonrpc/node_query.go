package jsonrpc

import (
	"github.com/emc-protocol/edge-matrix/application"
)

// RtcQuery is a query to filter rtc subjects
type NodeQuery struct {
	name     string
	nodeType string
	id       string
	version  string
}

// UnmarshalJSON decodes a json object
//
//	func (q *RtcQuery) UnmarshalJSON(data []byte) error {
//		var obj struct {
//			subjects     []string	`json:"subjects"`
//			applications []string	`json:"applications"`
//		}
//
//		err := json.Unmarshal(data, &obj)
//
//		if err != nil {
//			return err
//		}
//
//		if obj.subjects != nil {
//			// decode topics, either "" or ["", ""] or null
//			for _, item := range obj.subjects {
//
//			}
//		}
//
//		// decode topics
//		return nil
//	}
//
// Match returns whether the receipt includes topics for this filter
func (q *NodeQuery) Match(rm *application.Application) bool {
	if len(q.nodeType) > 0 {
		match := false
		if q.nodeType == rm.Tag {
			match = true
		}
		if !match {
			return false
		}
	}
	// check name
	if len(q.name) > 0 {
		match := false
		if rm.Name == q.name {
			match = true
		}

		if !match {
			return false
		}
	}

	if len(q.id) > 0 {
		match := false
		if rm.PeerID.String() == q.id {
			match = true
		}

		if !match {
			return false
		}
	}

	if len(q.version) > 0 {
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
