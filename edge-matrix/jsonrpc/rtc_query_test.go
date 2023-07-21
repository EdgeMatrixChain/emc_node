package jsonrpc

import (
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	app1 = "edge_chat"
	sub1 = types.StringToAddress("0x1234")
	sub2 = types.StringToAddress("0x0123")
)

func TestDecodeRtcQueryFromInterface(t *testing.T) {
	data := []byte(`{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "edge_subscribe",
  "params": [
    "rtc",
    {
      "applications": 
        "edge_chat"
      ,
      "subjects": 
        "0x0234"
     
    }
  ]
}`)

	query, err := decodeRtcQueryFromInterface(data)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "edge_chat", query.Application)

}
