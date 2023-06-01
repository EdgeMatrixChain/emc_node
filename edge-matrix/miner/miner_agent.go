package miner

import (
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/hashicorp/go-hclog"
	"math/big"
)

const (
	DEFAULT_MINER_CANISTER_ID = "nk6pr-3qaaa-aaaam-abnrq-cai"
)

var DEFAULT_IC_HOST = "https://ic0.app"

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

type myNodeResult struct {
	NodeInfos []NodeInfo
}

type NULL *uint8

type NodeInfo struct {
	//			"nodeID":     new(idl.Text),
	//			"owner":      new(idl.Principal),
	//			"wallet":     new(idl.Principal),
	//			"registered": new(idl.Nat),
	//			"nodeType":   new(idl.Nat),
	NodeID     string              `ic:"nodeID"`
	Owner      principal.Principal `ic:"owner"`
	Wallet     principal.Principal `ic:"wallet"`
	Registered big.Int             `ic:"registered"`
	NodeType   big.Int             `ic:"nodeType"`
}

type UpdateResult struct {
	// (variant {
	//	 Ok:nat;
	//	 Err:variant {StakeTooShort; NodeAlreadyExist; NodeNotExist; StakedBefore; UnknowType; StakeNotEnough; CallerNotAuthorized; DuplicatedValidation; CanNotUnstake; TokenTransferFailed; NoStakeFound; NotAValidator
	// }})
	Ok  big.Int `ic:"Ok"`
	Err EnumErr `ic:"Err"`
}

type EnumErr struct {
	// (variant {StakeTooShort; NodeAlreadyExist; NodeNotExist; StakedBefore; UnknowType; StakeNotEnough; CallerNotAuthorized; DuplicatedValidation; CanNotUnstake; TokenTransferFailed; NoStakeFound; NotAValidator)
	StakeTooShort        NULL `ic:"StakeTooShort"`
	NodeAlreadyExist     NULL `ic:"NodeAlreadyExist"`
	NodeNotExist         NULL `ic:"NodeNotExist"`
	StakedBefore         NULL `ic:"StakedBefore"`
	UnknowType           NULL `ic:"UnknowType"`
	StakeNotEnough       NULL `ic:"StakeNotEnough"`
	CallerNotAuthorized  NULL `ic:"CallerNotAuthorized"`
	DuplicatedValidation NULL `ic:"DuplicatedValidation"`
	CanNotUnstake        NULL `ic:"CanNotUnstake"`
	TokenTransferFailed  NULL `ic:"TokenTransferFailed"`
	NoStakeFound         NULL `ic:"NoStakeFound"`
	NotAValidator        NULL `ic:"NotAValidator"`
	//To formulate a enum struct
	Index string `ic:"EnumIndex"`
}

// call miner canister's myNode method(query)
func (m *MinerAgent) MyNode(nodeId string) (string, string, string, int64, int64, error) {
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
		return "", "", "", -1, -1, err
	}
	m.logger.Debug("MyNode", "types", types, "result", result)
	//fmt.Println("types:", types)
	//ouput-> types: [interface vec interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}]
	//fmt.Println("result:", result)
	//ouput-> result: [[map[3054210041:[57 248 44 186 16 145 100 93 182 123 49 153 147 45 214 150 45 224 237 216 84 142 130 10 172 82 193 48 2] 4104166786:1685512050331435992 4135997916:0 656559709:16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC 947296307:[245 50 153 79 90 148 3 179 181 210 38 205 150 98 51 71 55 221 150 24 248 186 191 134 143 61 52 87 2]]]]
	if result != nil && len(result) > 0 {
		recordVec := result[0].([]interface{})
		if len(recordVec) > 0 {
			record := recordVec[0]
			myNode := NodeInfo{}
			utils.Decode(&myNode, record.(map[string]interface{}))
			return myNode.NodeID, myNode.Owner.Encode(), myNode.Wallet.Encode(), myNode.Registered.Int64(), myNode.NodeType.Int64(), nil
		}
	}
	return "", "", "", -1, -1, nil
}

// call miner canister's myCurrentEPower method(query)
func (m *MinerAgent) MyCurrentEPower(nodeId string) (uint64, uint64, error) {
	methodName := "myCurrentEPower"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	argValue := []interface{}{
		nodeId,
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("MyCurrentEPower", "argType", argType, "argValue", argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		return 0, 0, err
	}
	m.logger.Debug("MyCurrentEPower", "types", types, "result", result)
	fmt.Println("types:", types)
	//ouput-> types: [interface vec interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}]
	fmt.Println("result:", result)
	//ouput-> result: [[map[3054210041:[57 248 44 186 16 145 100 93 182 123 49 153 147 45 214 150 45 224 237 216 84 142 130 10 172 82 193 48 2] 4104166786:1685512050331435992 4135997916:0 656559709:16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC 947296307:[245 50 153 79 90 148 3 179 181 210 38 205 150 98 51 71 55 221 150 24 248 186 191 134 143 61 52 87 2]]]]
	if result != nil && len(result) >= 2 {
		return result[0].(*big.Int).Uint64(), result[0].(*big.Int).Uint64(), nil
	}
	return 0, 0, err
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
	if len(result) < 1 {
		return errors.New("result is empty")
	}

	if result != nil && len(result) > 0 {
		respVariant := result[0].(map[string]interface{})
		updateResult := UpdateResult{}
		utils.Decode(&updateResult, respVariant)
		if updateResult.Ok.Int64() > 0 {
			m.logger.Info("RegisterNode ok")
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("RegisterNode fail")
}

// call miner canister's unRegisterNode method(update)
func (m *MinerAgent) UnRegisterNode(nodeId string) error {
	methodName := "unregisterNode"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	argValue := []interface{}{
		nodeId,
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("UnRegisterNode", "argType", argType, "argValue", argValue)

	types, result, err := m.agent.Update(m.canister, methodName, arg, 30)
	if err != nil {
		return err
	}
	m.logger.Debug("UnRegisterNode", "types", types[0].String(), "result", result)
	if len(result) < 1 {
		return errors.New("result is empty")
	}

	if result != nil && len(result) > 0 {
		respVariant := result[0].(map[string]interface{})
		updateResult := UpdateResult{}
		utils.Decode(&updateResult, respVariant)
		if updateResult.Ok.Int64() > 0 {
			m.logger.Info("UnRegisterNode ok")
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("UnRegisterNode fail")
}

// call miner canister's submitValidation method(update)
func (m *MinerAgent) SubmitValidation(validationTicket int64, validator string, power int64, targetNodeID string) error {
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
				})))

	argValue := []interface{}{
		[]interface{}{
			map[string]interface{}{
				"validationTicket": big.NewInt(validationTicket),
				"validator":        p,
				"power":            big.NewInt(power),
				"targetNodeID":     targetNodeID,
			}},
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("SubmitValidation", "argType", argType, "argValue", argValue)

	types, result, err := m.agent.Update(m.canister, methodName, arg, 30)
	if err != nil {
		return err
	}
	m.logger.Debug("SubmitValidation", "types", types[0].String(), "result", result)
	//fmt.Println("result:", result)
	// (variant {Ok=0}) -> result: [map[17724:0 EnumIndex:17724]]
	// (variant {Err=variant {NotAValidator}}) -> result: [map[3456837:map[3734858244:<nil> EnumIndex:3734858244] EnumIndex:3456837]]
	if len(result) < 1 {
		return errors.New("result is empty")
	}

	if result != nil && len(result) > 0 {
		respVariant := result[0].(map[string]interface{})
		updateResult := UpdateResult{}
		utils.Decode(&updateResult, respVariant)
		if updateResult.Ok.Int64() > 0 {
			m.logger.Info("SubmitValidation ok")
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("SubmitValidation fail")
}
