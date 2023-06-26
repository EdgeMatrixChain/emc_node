package miner

import (
	"errors"
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

type ValidorNodeInfo struct {
	//{nodeID:text; owner:principal; wallet:principal; registered:int; nodeType:nat}}
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
func (m *MinerAgent) MyCurrentEPower(nodeId string) (uint64, float32, error) {
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
	//ouput-> types: [interface vec interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}]
	//ouput-> result: [[map[3054210041:[57 248 44 186 16 145 100 93 182 123 49 153 147 45 214 150 45 224 237 216 84 142 130 10 172 82 193 48 2] 4104166786:1685512050331435992 4135997916:0 656559709:16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC 947296307:[245 50 153 79 90 148 3 179 181 210 38 205 150 98 51 71 55 221 150 24 248 186 191 134 143 61 52 87 2]]]]
	if result != nil && len(result) >= 2 {
		total, _ := result[1].(*big.Float).Float32()
		return result[0].(*big.Int).Uint64(), total, nil
	}
	return 0, 0, err
}

func (m *MinerAgent) MyStack(nodeId string) (uint64, uint64, uint64, error) {
	methodName := "myStake"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	argValue := []interface{}{
		nodeId,
	}
	arg, _ := idl.Encode(argType, argValue)
	m.logger.Debug("MyStack", "argType", argType, "argValue", argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		return 0, 0, 10000, err
	}
	m.logger.Debug("MyStack", "types", types, "result", result)
	//ouput-> types: [interface vec interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}]
	//ouput-> result: [[map[3054210041:[57 248 44 186 16 145 100 93 182 123 49 153 147 45 214 150 45 224 237 216 84 142 130 10 172 82 193 48 2] 4104166786:1685512050331435992 4135997916:0 656559709:16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC 947296307:[245 50 153 79 90 148 3 179 181 210 38 205 150 98 51 71 55 221 150 24 248 186 191 134 143 61 52 87 2]]]]
	if result != nil && len(result) >= 3 {
		return result[0].(*big.Int).Uint64(), result[1].(*big.Int).Uint64(), result[2].(*big.Int).Uint64(), nil
	}
	return 0, 0, 10000, err
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
			m.logger.Info("\n------------------------------------------\nSubmitValidation ok", "targetNodeID", targetNodeID)
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("SubmitValidation fail")
}

func (m *MinerAgent) SubmitValidationVec(vecValue []interface{}) error {
	methodName := "submitValidation"
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

	argValue := []interface{}{vecValue}
	arg, _ := idl.Encode(argType, argValue)
	//fmt.Println(argType, "   ", argValue)
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
			return nil
		} else {
			return errors.New(updateResult.Err.Index)
		}
	}
	return errors.New("SubmitValidation fail")
}

func (m *MinerAgent) ListValidatorsNodeId() ([]string, error) {
	nodeList := make([]string, 0)
	methodName := "listValidatorNodes"
	var argType []idl.Type
	var argValue []interface{}
	arg, _ := idl.Encode(argType, argValue)

	types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	if err != nil {
		m.logger.Error(err.Error())
		return nodeList, err
	}
	//fmt.Println("types:", types[0].String(), "result:", result)
	//types: interface vec interface record {0:text; 1:interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}} result: [[map[0:16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y 947296307:[36 168 117 208 89 72 247 191 251 207 187 65 155 215 201 21 37 102 213 51 157 96 202 221 113 185 26 84 2]]] map[0:16Uiu2HAmEoDReK7pKygYYYFgJ8uuXS8oWsYFWiiEbCSF9HjYcih2 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmEoDReK7pKygYYYFgJ8uuXS8oWsYFWiiEbCSF9HjYcih2 947296307:[4 147 74 93 122 138 141 1 188 183 152 85 46 153 244 27 237 124 74 239 66 112 148 29 10 21 174 46 2]]] map[0:16Uiu2HAkyPw8SEeDpErEwcEZ2QtXzPq5KQf4woWybsr7KN6VH7yX 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685612201271035791 4135997916:1 656559709:16Uiu2HAkyPw8SEeDpErEwcEZ2QtXzPq5KQf4woWybsr7KN6VH7yX 947296307:[137 93 57 21 145 14 171 17 20 8 27 60 176 216 139 147 110 198 119 130 170 1 175 38 61 19 108 211 2]]] map[0:16Uiu2HAmPfFVHNnYKdDQywJXnzbgM1MdAi6P1MsCkxN7Hr6VaiYa 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmPfFVHNnYKdDQywJXnzbgM1MdAi6P1MsCkxN7Hr6VaiYa 947296307:[42 174 237 43 160 163 242 37 48 118 200 166 112 189 56 255 19 59 90 30 53 188 133 153 30 253 236 40 2]]] map[0:16Uiu2HAmTPfBgUkQ4V8qaBvTaJp54Cm32TWGvYZaxcuPxoaSbZAS 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmTPfBgUkQ4V8qaBvTaJp54Cm32TWGvYZaxcuPxoaSbZAS 947296307:[102 242 43 171 69 192 233 201 175 42 233 72 161 61 119 192 214 32 225 200 111 96 26 111 57 76 242 248 2]]] map[0:16Uiu2HAmGpKZdnpaaYgKTZqagLVJcnMphdeqaHtKBaFFkb5MYRUy 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmGpKZdnpaaYgKTZqagLVJcnMphdeqaHtKBaFFkb5MYRUy 947296307:[68 196 3 132 35 19 140 71 180 180 247 198 72 205 39 164 52 180 106 88 60 107 211 100 131 101 160 103 2]]] map[0:16Uiu2HAm7BqtmjH7JECa5Y4iNgiZXuet3HqXYZbeXNN9XiQwgSbf 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605847289287331 4135997916:1 656559709:16Uiu2HAm7BqtmjH7JECa5Y4iNgiZXuet3HqXYZbeXNN9XiQwgSbf 947296307:[46 113 59 128 101 160 251 203 212 119 205 85 133 9 206 253 168 134 11 30 174 45 49 197 147 0 246 137 2]]]]]
	m.logger.Debug("ListValidators", "types", types[0].String(), "result", result)
	//fmt.Println("len(result):", len(result))
	if result != nil && len(result) > 0 {
		vec := result[0].([]interface{})
		for _, recordVec := range vec {
			record := recordVec.(map[string]interface{})
			nodeList = append(nodeList, record["0"].(string))
		}
		return nodeList, nil
	}
	return nil, nil
}
