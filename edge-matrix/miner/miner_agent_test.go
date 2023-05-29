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
	privKeyBytes, err := hex.DecodeHex(privKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, privKey)
	//minerAgent := NewMinerAgentWithICKey(hclog.NewNullLogger(), config.IcHost, hex.EncodeToString(decodedPrivKey.Seed()), config.MinerCanister)

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	//minerAgent := NewMinerAgentWithICKey(hclog.NewNullLogger(), "http://127.0.0.1:8081", privKey, DEFAULT_MINER_CANISTER_ID)
	err = minerAgent.RegisterNode(
		NodeTypeRouter,
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
		"rlqvd-pzz7a-wluee-rmro3-m6zrt-gjs3v-uwfxq-o3wcu-r2bav-lcsye-yae")
	if err != nil {
		t.Log(err)
	}
}
