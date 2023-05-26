package rtc

import (
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/types"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRLPMarshall_And_Unmarshall_RtcMsg(t *testing.T) {
	addrFrom := types.StringToAddress("12")
	msg := &RtcMsg{
		Subject:     "2",
		Application: "3",
		Content:     "4",
		V:           big.NewInt(25),
		S:           big.NewInt(26),
		R:           big.NewInt(27),
		From:        addrFrom,
		Type:        SubjectMsg,
		To:          types.ZeroAddress,
	}
	unmarshalledMsg := new(RtcMsg)
	marshaledRlp := msg.MarshalRLP()

	t.Log("marshaledRlp:" + hex.EncodeToHex(marshaledRlp))
	if err := unmarshalledMsg.UnmarshalRLP(marshaledRlp); err != nil {
		t.Fatal(err)
	}

	unmarshalledMsg.ComputeHash()
	t.Log("unmarshalledContent" + unmarshalledMsg.Content)
	msg.Hash = unmarshalledMsg.Hash
	assert.Equal(t, msg, unmarshalledMsg, "[ERROR] Unmarshalled rtcMsg not equal to base rtcMsg")
}
