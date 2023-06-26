package miner

import (
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/identity"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/hashicorp/go-hclog"
	"math/big"
	"testing"
	"time"
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

	count, e, err := minerAgent.MyCurrentEPower(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
	t.Log("count:", count)
	t.Log("e:", e)
}

func Test_MyStack(t *testing.T) {
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

	_, _, multiple, err := minerAgent.MyStack(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
	t.Log("multiple:", multiple)
	t.Log("rate:", float32(multiple)/10000.0)
}

func Test_sub(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7}
	arr0 := arr[:2]
	arr1 := arr[2:7]
	t.Log(arr0)
	t.Log(arr1)
}

func Test_SubmitValidationVec(t *testing.T) {
	validatorPrivKey := privKey
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, validatorPrivKey)
	privKeyBytes, err := hex.DecodeHex(validatorPrivKey)

	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	//vecValues := []interface{}{
	//	map[string]interface{}{
	//		"validationTicket": big.NewInt(1000),
	//		"validator":        p,
	//		"power":            big.NewInt(150000),
	//		"targetNodeID":     "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
	//	},
	//	map[string]interface{}{
	//		"validationTicket": big.NewInt(1200),
	//		"validator":        p,
	//		"power":            big.NewInt(150000),
	//		"targetNodeID":     "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
	//	}}
	batchSize := 10
	batchSubmitData := make([]*proof.PocSubmitData, batchSize)

	taskCount := 0
	success := 0
	for {
		time.Sleep(20 * time.Millisecond)
		batchSubmitData[taskCount] = &proof.PocSubmitData{
			Validator:        p.Encode(),
			Power:            150000,
			TargetNodeID:     "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
			ValidationTicket: int64(1000 + taskCount),
		}

		taskCount += 1

		if taskCount < batchSize {
			continue
		}

		taskCount = 0
		vecValues := make([]interface{}, len(batchSubmitData))
		for i, pocSubmitData := range batchSubmitData {
			p, err := principal.Decode(pocSubmitData.Validator)
			if err != nil {
				t.Error("principal.Decode", "err", err)
				continue
			}

			vecValues[i] = map[string]interface{}{
				"validationTicket": big.NewInt(pocSubmitData.ValidationTicket),
				"validator":        p,
				"power":            big.NewInt(pocSubmitData.Power),
				"targetNodeID":     pocSubmitData.TargetNodeID,
			}
		}
		submitToIc(t, minerAgent, vecValues)
		success += 1
		if success > 1 {
			return
		}
	}
}
func submitToIc(t *testing.T, minerAgent *MinerAgent, vecValues []interface{}) {
	// submit proof result to IC canister
	t.Log("\n------------------------------------------\nSubmitValidation", "posting...", len(vecValues))
	err := minerAgent.SubmitValidationVec(vecValues)
	if err != nil {
		t.Error("\n------------------------------------------\nSubmitValidation", "err", err)
	} else {
		t.Log("\n------------------------------------------\nSubmitValidation", "success", len(vecValues))
	}
}

func Test_SubmitValidation(t *testing.T) {
	validatorPrivKey := privKey
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, validatorPrivKey)
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
