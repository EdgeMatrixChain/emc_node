package miner

import (
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/identity"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/hashicorp/go-hclog"
	"testing"
)

const (
	privKey = "ed1cb741ef10f2e353c6c395f0b91270e762e55cf30d03ab6fa340ff306fb9d9"
)

func Test_RegisterNode(t *testing.T) {
	icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	//minerAgent := NewMinerAgentWithICKey(hclog.NewNullLogger(), "http://127.0.0.1:8081", privKey, DEFAULT_MINER_CANISTER_ID)

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.RegisterNode(
		NodeTypeRouter,
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
		"rlqvd-pzz7a-wluee-rmro3-m6zrt-gjs3v-uwfxq-o3wcu-r2bav-lcsye-yae")
	if err != nil {
		t.Log(err)
	}
}

func Test_MyNode(t *testing.T) {
	icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)
	wp, nodeType, err := minerAgent.MyNode(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
	t.Log("wallet:", wp)
	t.Log("nodeType:", nodeType)
}

func Test_SubmitValidation(t *testing.T) {
	icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)
	err = minerAgent.SubmitValidation(
		1000,
		p.Encode(),
		150000,
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
		NodeTypeRouter,
	)
	if err != nil {
		t.Log(err)
	}
}
