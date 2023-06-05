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
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.RegisterNode(
		NodeTypeRouter,
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
		"rlqvd-pzz7a-wluee-rmro3-m6zrt-gjs3v-uwfxq-o3wcu-r2bav-lcsye-yae")
	if err != nil {
		t.Log(err)
	}
}

func Test_UnRegisterNode(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.UnRegisterNode(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
}

func Test_MyNode(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	nodeId, nodeIdentity, wp, registered, nodeType, err := minerAgent.MyNode(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
	t.Log("nodeId:", nodeId)
	t.Log("nodeIdentity:", nodeIdentity)
	t.Log("wallet:", wp)
	t.Log("registered:", registered > 0)
	t.Log("nodeType:", nodeType)
}

func Test_MyCurrentEPower(t *testing.T) {
	icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	count, e, err := minerAgent.MyCurrentEPower(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
	t.Log("count:", count)
	t.Log("e:", e)
}

func Test_SubmitValidation(t *testing.T) {
	validatorPrivKey := "8031dda21dd9a138a93c9c60ac866608c6f0f8d1a8a79ffbd7faf59faeb2b1d7"
	//validatorPrivKey := privKey
	icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, validatorPrivKey)
	privKeyBytes, err := hex.DecodeHex(validatorPrivKey)

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
	)
	if err != nil {
		t.Log(err)
	}
}
func Test_listValidators(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	icAgent := agent.NewWithHost("https://ic0.app", false, privKey)
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	nodeList, err := minerAgent.ListValidatorsNodeId()
	if err != nil {
		t.Log(err)
	}
	t.Log(len(nodeList))
	for _, validatorNodeID := range nodeList {
		t.Log(validatorNodeID)
	}
}
