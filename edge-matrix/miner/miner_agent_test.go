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
	TestPrivKey = "ed1cb741ef10f2e353c6c395f0b91270e762e55cf30d03ab6fa340ff306fb9d9"
)

func Test_RegisterNode(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.RegisterComputingNode(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
		"c6wnt-id2x5-dz525-mcrjd-tkor7-jkjjl-gaqbz-f5uia-obawe-azbxs-qae")
	if err != nil {
		t.Log(err)
	}
}

func Test_AddRouter(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.AddRouter(
		"c6wnt-id2x5-dz525-mcrjd-tkor7-jkjjl-gaqbz-f5uia-obawe-azbxs-qae")
	if err != nil {
		t.Log(err)
	}
}

type nodeData struct {
	nodeId         string
	privKey        string
	nodeType       NodeType
	minerPrincipal string
}

func Test_UnRegisterNode(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.UnRegisterComputingNode(
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		t.Log(err)
	}
}

func Test_MyNode(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	//minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, "be2us-64aaa-aaaaa-qaabq-cai")
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
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
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
	if count > 0 {
		t.Log("average:", e/float32(count))

	}
}

func Test_MyStack(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
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

func Test_SubmitValidationVec(t *testing.T) {
	validatorPrivKey := TestPrivKey
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
	//		"power":            big.NewInt(150000),
	//		"targetNodeID":     "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
	//	},
	//	map[string]interface{}{
	//		"validationTicket": big.NewInt(1200),
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
			vecValues[i] = map[string]interface{}{
				"validationTicket": big.NewInt(pocSubmitData.ValidationTicket),
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

func Test_SubmitComputingValidation(t *testing.T) {
	validatorPrivKey := TestPrivKey
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, validatorPrivKey)
	privKeyBytes, err := hex.DecodeHex(validatorPrivKey)

	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	minerAgent := NewMinerAgent(hclog.NewNullLogger(), icAgent, DEFAULT_MINER_CANISTER_ID)

	err = minerAgent.SubmitComputingValidation(
		1000,
		150000,
		"16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC",
	)
	if err != nil {
		t.Log(err)
	}
}

func Test_listValidators(t *testing.T) {
	//icAgent := agent.NewWithHost("http://127.0.0.1:8081", false, TestPrivKey)
	icAgent := agent.NewWithHost("https://ic0.app", false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
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
