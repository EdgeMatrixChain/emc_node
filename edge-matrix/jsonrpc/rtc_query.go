package jsonrpc

import (
	"encoding/json"
	"github.com/emc-protocol/edge-matrix/rtc"
	"github.com/emc-protocol/edge-matrix/types"
	"math/big"
)

// RtcQuery is a query to filter rtc subjects
type RtcQuery struct {
	Subject     string
	Application string
	From        string
}

func decodeRtcQueryFromInterface(i interface{}) (*RtcQuery, error) {
	// once the rtc filter is decoded as map[string]interface we cannot use unmarshal json
	raw, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var obj struct {
		Subject     string `json:"subject"`
		Application string `json:"application"`
		Content     string `json:"content"`
		V           string `json:"V"`
		R           string `json:"R"`
		S           string `json:"S"`
	}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	bigV := new(big.Int)
	bigV.SetString(obj.V, 0)
	bigR := new(big.Int)
	bigR.SetString(obj.R, 0)
	bigS := new(big.Int)
	bigS.SetString(obj.S, 0)
	query := &RtcQuery{
		Subject:     obj.Subject,
		Application: obj.Application,
	}
	return query, nil
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
func (q *RtcQuery) Match(rm *rtc.RtcMsg) bool {
	// check addresses
	// TODO if has To filed in msg
	if rm.To != types.ZeroAddress && rm.To != types.StringToAddress(q.From) {
		return false
	}

	if len(q.Application) > 0 {
		match := false
		if q.Application == rm.Application {
			match = true
		}
		if !match {
			return false
		}
	}
	// check subjects
	if len(q.Subject) > 0 {
		match := false
		if rm.Subject == q.Subject {
			match = true
		}

		if !match {
			return false
		}
	}

	return true
}
