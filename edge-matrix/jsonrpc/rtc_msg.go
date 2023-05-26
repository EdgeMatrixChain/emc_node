package jsonrpc

import (
	"encoding/json"
	"github.com/emc-protocol/edge-matrix/rtc"
)

type RtcMsg struct {
	//Nonce       uint64
	Subject     string
	Application string
	Content     string
	To          string

	V string
	R string
	S string

	Type rtc.RtcType
}

func DecodeRtcMsgFromInterface(i interface{}) (*RtcMsg, error) {
	raw, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	obj := &RtcMsg{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}
