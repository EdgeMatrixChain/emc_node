package miner

import (
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/hashicorp/go-hclog"
	"math/big"
)

const (
	DEFAULT_MINER_CANISTER_ID = "bw4dl-smaaa-aaaaa-qaacq-cai"
)

type MinerAgent struct {
	logger hclog.Logger

	agent    *agent.Agent
	canister string
}

type NodeType int64

const (
	NodeTypeRouter    NodeType = 0
	NodeTypeValidator NodeType = 1
	NodeTypeComputing NodeType = 2
)

func NewMinerAgentWithICKey(logger hclog.Logger, icHost string, icPrivKey string, canister string) *MinerAgent {
	return &MinerAgent{
		logger:   logger,
		agent:    agent.NewWithHost(icHost, false, icPrivKey),
		canister: canister,
	}
}

func NewMinerAgent(logger hclog.Logger, agent *agent.Agent, canister string) *MinerAgent {
	return &MinerAgent{
		logger:   logger,
		agent:    agent,
		canister: canister,
	}
}

func NewMinerAgentWithCanister(logger hclog.Logger, icHost string, icPrivKey string, minerCanister string) *MinerAgent {
	return &MinerAgent{
		logger:   logger,
		agent:    agent.NewWithHost(icHost, false, icPrivKey),
		canister: minerCanister,
	}
}

func (m *MinerAgent) RegisterNode(nodeType NodeType, nodeId string, minerPrincipal string) error {
	methodName := "registerNode"

	p, err := principal.Decode(minerPrincipal)
	if err != nil {
		return err
	}
	var argType []idl.Type
	argType = append(argType, new(idl.Nat))
	argType = append(argType, new(idl.Text))
	argType = append(argType, new(idl.Principal))

	argValue := []interface{}{
		big.NewInt(int64(nodeType)),
		nodeId,
		p}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("RegisterNode", "argType", argType, "argValue", argValue)
	fmt.Println("RegisterNode", "argType", argType, "argValue", argValue)

	types, result, err := m.agent.Update(m.canister, methodName, arg, 30)
	if err != nil {
		return err
	}
	m.logger.Debug("RegisterNode", "types", types[0].String(), "result", result)
	// (variant {Ok=0}) -> [map[17724:0 EnumIndex:17724]]
	// (variant {Err=variant {NodeAlreadyExist}})-> [map[3456837:map[440058177:<nil> EnumIndex:440058177] EnumIndex:3456837]]
	if len(result) < 1 {
		return errors.New("result is empty")
	}

	resultMap := result[0].(map[string]interface{})
	valueKeyName := ""
	for k, v := range resultMap {
		if k == "EnumIndex" {
			valueKeyName = v.(string)
		}
	}
	if valueKeyName != "" {
		value := resultMap[valueKeyName]
		if mapValue, ok := value.(map[string]interface{}); ok {
			m.logger.Error("RegisterNode", "err", mapValue)
			return errors.New("register fail")
		} else {
			m.logger.Info("RegisterNode", "ok", value)
		}

	}
	return nil
}
