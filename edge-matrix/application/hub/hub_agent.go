package hub

import (
	"errors"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
	"github.com/hashicorp/go-hclog"
	"math/big"
)

const (
	DEFAULT_HUB_CANISTER_ID = "57ab7-fiaaa-aaaam-abr2q-cai"
)

var DEFAULT_IC_HOST = "https://ic0.app"

type NULL *uint8

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
	UnknownError        NULL `ic:"UnknownError"`
	ModelAlreadyExist   NULL `ic:"ModelAlreadyExist"`
	CallerNotAuthorized NULL `ic:"CallerNotAuthorized"`
	ModelNotExist       NULL `ic:"ModelNotExist"`
	//To formulate a enum struct
	Index string `ic:"EnumIndex"`
}

type HubAgent struct {
	logger hclog.Logger

	agent    *agent.Agent
	canister string
}

func NewHubAgentWithICKey(logger hclog.Logger, icHost string, icPrivKey string, canister string) *HubAgent {
	return &HubAgent{
		logger:   logger,
		agent:    agent.NewWithHost(icHost, false, icPrivKey),
		canister: canister,
	}
}

func NewHubAgent(logger hclog.Logger, agent *agent.Agent) *HubAgent {
	return &HubAgent{
		logger:   logger,
		agent:    agent,
		canister: DEFAULT_HUB_CANISTER_ID,
	}
}

func NewHubAgentWithCanister(logger hclog.Logger, icHost string, icPrivKey string, minerCanister string) *HubAgent {
	return &HubAgent{
		logger:   logger,
		agent:    agent.NewWithHost(icHost, false, icPrivKey),
		canister: minerCanister,
	}
}

func (m *HubAgent) GetIdentity() string {
	return m.agent.GetIdentity()
}

type HubModelInfo struct {
	//{nodeID:text; owner:principal; wallet:principal; registered:int; nodeType:nat}}
	ModelHash string `ic:"0"`
	ModelType string `ic:"1"`
}

func (m *HubAgent) AddModel(modelType string, modelHash string) error {
	methodName := "addModel"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))
	argType = append(argType, new(idl.Text))

	argValue := []interface{}{
		modelType,
		modelHash,
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("AddModel", "argType", argType, "argValue", argValue)

	types, result, err := m.agent.Update(m.canister, methodName, arg, 30)
	if err != nil {
		return err
	}
	m.logger.Debug("AddModel", "types", types[0].String(), "result", result)
	if len(result) < 1 {
		return errors.New("result is empty")
	}

	if result != nil && len(result) > 0 {
		respVariant := result[0].(map[string]interface{})
		updateResult := UpdateResult{}
		utils.Decode(&updateResult, respVariant)
		if updateResult.Ok.Int64() > 0 {
			m.logger.Info("AddModel ok")
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("AddModel fail")
}

func (m *HubAgent) RemoveModel(modelHash string) error {
	methodName := "removeModel"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	argValue := []interface{}{
		modelHash,
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("RemoveModel", "argType", argType, "argValue", argValue)

	types, result, err := m.agent.Update(m.canister, methodName, arg, 30)
	if err != nil {
		return err
	}
	m.logger.Debug("RemoveModel", "types", types[0].String(), "result", result)
	if len(result) < 1 {
		return errors.New("result is empty")
	}

	if result != nil && len(result) > 0 {
		respVariant := result[0].(map[string]interface{})
		updateResult := UpdateResult{}
		utils.Decode(&updateResult, respVariant)
		if updateResult.Ok.Int64() > 0 {
			m.logger.Info("RemoveModel ok")
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("RemoveModel fail")
}

func (m *HubAgent) IsModelListed(modelHash string) (bool, error) {
	methodName := "isModelListed"
	var argType []idl.Type
	argType = append(argType, new(idl.Text))
	argValue := []interface{}{
		modelHash,
	}
	arg, _ := idl.Encode(argType, argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		m.logger.Error(err.Error())
		return false, err
	}
	//fmt.Println("types:", types[0].String(), "result:", result)
	// types: bool result: [true]
	m.logger.Debug("IsModelListed", "types", types[0].String(), "result", result)
	if result != nil && len(result) > 0 {
		isValid := result[0].(bool)
		return isValid, nil
	}
	return false, nil
}

func (m *HubAgent) ListModels() ([]string, error) {
	modelList := make([]string, 0)
	methodName := "listModels"
	var argType []idl.Type
	var argValue []interface{}
	arg, _ := idl.Encode(argType, argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		m.logger.Error(err.Error())
		return modelList, err
	}
	//fmt.Println("types:", types[0].String(), "result:", result)
	//listModels types interface vec interface record {0:text; 1:text} result [[map[0:hash0000003 1:openai] map[0:hash0000002 1:stable] map[0:hash0000001 1:stable]]]
	m.logger.Debug("listModels", "types", types[0].String(), "result", result)
	if result != nil && len(result) > 0 {
		vec := result[0].([]interface{})
		for _, recordVec := range vec {
			record := recordVec.(map[string]interface{})

			modelInfo := HubModelInfo{}
			utils.Decode(&modelInfo, record)

			modelList = append(modelList, modelInfo.ModelHash)
		}
		return modelList, nil
	}
	return nil, nil
}

func (m *HubAgent) ListModelsByeType(modelType string) ([]string, error) {
	modelList := make([]string, 0)
	methodName := "listModelsByeType"
	var argType []idl.Type
	argType = append(argType, new(idl.Text))
	argValue := []interface{}{
		modelType,
	}
	arg, _ := idl.Encode(argType, argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		m.logger.Error(err.Error())
		return modelList, err
	}
	//fmt.Println("types:", types[0].String(), "result:", result)
	//listModels types interface vec interface record {0:text; 1:text} result [[map[0:hash0000003 1:openai] map[0:hash0000002 1:stable] map[0:hash0000001 1:stable]]]
	m.logger.Debug("listModels", "types", types[0].String(), "result", result)
	if result != nil && len(result) > 0 {
		vec := result[0].([]interface{})
		for _, recordVec := range vec {
			record := recordVec.(map[string]interface{})

			modelInfo := HubModelInfo{}
			utils.Decode(&modelInfo, record)

			modelList = append(modelList, modelInfo.ModelHash)
		}
		return modelList, nil
	}
	return nil, nil
}
