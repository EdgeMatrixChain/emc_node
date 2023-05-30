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

func (m *MinerAgent) GetIdentity() string {
	return m.agent.GetIdentity()
}

// call miner canister's myNode method(query)
func (m *MinerAgent) MyNode(nodeId string) (string, int64, error) {
	methodName := "myNode"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	argValue := []interface{}{
		nodeId,
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("MyNode", "argType", argType, "argValue", argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		return "", -1, err
	}
	m.logger.Debug("MyNode", "types", types, "result", result)
	//fmt.Println("result:", result)
	// [[map[0:[57 248 44 186 16 145 100 93 182 123 49 153 147 45 214 150 45 224 237 216 84 142 130 10 172 82 193 48 2] 1:0]]]
	if result != nil && len(result) > 0 {
		vec := result[0].([]interface{})
		if len(vec) > 0 {
			record := vec[0].(map[string]interface{})
			return principal.New(record["0"].([]byte)).Encode(), record["1"].(*big.Int).Int64(), nil
		}
	}
	return "", -1, nil
}

// call miner canister's registerNode method(update)
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

// call miner canister's submitValidation method(update)
func (m *MinerAgent) SubmitValidation(validationTicket int64, validator string, power int64, targetNodeID string, nodeType NodeType) error {
	methodName := "submitValidation"

	p, err := principal.Decode(validator)
	if err != nil {
		return err
	}
	var argType []idl.Type
	argType = append(
		argType,
		idl.NewVec(
			idl.NewRec(
				map[string]idl.Type{
					"validationTicket": new(idl.Nat),
					"validator":        new(idl.Principal),
					"power":            new(idl.Nat),
					"targetNodeID":     new(idl.Text),
					"nodeType":         new(idl.Nat),
				})))

	argValue := []interface{}{
		[]interface{}{
			map[string]interface{}{
				"validationTicket": big.NewInt(validationTicket),
				"validator":        p,
				"power":            big.NewInt(power),
				"targetNodeID":     targetNodeID,
				"nodeType":         big.NewInt(int64(nodeType)),
			}},
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("SubmitValidation", "argType", argType, "argValue", argValue)

	types, result, err := m.agent.Update(m.canister, methodName, arg, 30)
	if err != nil {
		return err
	}
	m.logger.Debug("SubmitValidation", "types", types[0].String(), "result", result)
	fmt.Println("result:", result)
	// (variant {Ok=0}) -> result: [map[17724:0 EnumIndex:17724]]
	// (variant {Err=variant {NotAValidator}}) -> result: [map[3456837:map[3734858244:<nil> EnumIndex:3734858244] EnumIndex:3456837]]
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
			m.logger.Error("SubmitValidation", "err", mapValue)
			return errors.New("submit fail")
		} else {
			m.logger.Info("SubmitValidation", "ok", value)
		}

	}
	return nil
}
