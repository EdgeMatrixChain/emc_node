package miner

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/hashicorp/go-hclog"
	"math/rand"
)

var DEFAULT_HUB_HOST = "https://api.edgematrix.pro"

type MinerHubAgent struct {
	logger         hclog.Logger
	httpClient     *rpc.FastHttpClient
	secretsManager secrets.SecretsManager
}

type NodeType int64

const (
	NodeTypeRouter    NodeType = 0
	NodeTypeValidator NodeType = 1
	NodeTypeComputing NodeType = 2
)

type NULL *uint8

type NodeInfo struct {
	NodeID    string `json:"nodeId"`
	Principal string `json:"principal"`
	PublicKey string `json:"publicKey"`
	Status    int    `json:"status"`
	NodeType  string `json:"nodeType"`
}

type EPower struct {
	NodeID string  `json:"nodeId"`
	Round  uint64  `json:"round"`
	Power  float32 `json:"power"`
}

//type ValidorNodeInfo struct {
//	//{nodeID:text; owner:principal; wallet:principal; registered:int; nodeType:nat}}
//	NodeID     string  `json:"nodeID"`
//	Owner      string  `json:"owner"`
//	Wallet     string  `json:"wallet"`
//	Registered big.Int `json:"registered"`
//	NodeType   big.Int `json:"nodeType"`
//}

func NewMinerHubAgent(logger hclog.Logger, secretsManager secrets.SecretsManager) *MinerHubAgent {
	return &MinerHubAgent{
		logger:         logger,
		httpClient:     rpc.NewDefaultHttpClient(),
		secretsManager: secretsManager,
	}
}

// call EMCHub's query api
func (m *MinerHubAgent) MyCurrentEPower(nodeId string) (uint64, float32, error) {
	respBytes, err := m.httpClient.SendGetRequest(DEFAULT_HUB_HOST + "/api/v1/nodesign/query?nodeId=" + nodeId)
	if err != nil {
		return 0, 0, errors.New("Query EPower fail: " + err.Error())
	}

	if len(respBytes) > 0 {
		m.logger.Debug("Query EPower", "resp", string(respBytes))

		var response struct {
			Result int    `json:"_result"`
			Desc   string `json:"_desc"`
			Data   EPower `json:"data"`
		}

		err := json.Unmarshal(respBytes, &response)
		if err != nil {
			return 0, 0, errors.New("Query EPower fail: " + err.Error())
		}

		if response.Result == 0 {
			// return myNode.NodeID, myNode.Owner.Encode(), myNode.Wallet.Encode(), myNode.Registered.Int64(), myNode.NodeType.Int64(), nil
			return response.Data.Round, response.Data.Power, nil
		} else {
			return 0, 0, errors.New("Query EPower fail: " + err.Error())
		}
	}
	return 0, 0, errors.New("Query EPower fail")
}

// call EMCHub's query api
func (m *MinerHubAgent) MyNode(nodeId string) (string, string, string, int64, string, error) {
	respBytes, err := m.httpClient.SendGetRequest(DEFAULT_HUB_HOST + "/api/v1/nodesign/query?nodeId=" + nodeId)
	if err != nil {
		return "", "", "", -1, "", errors.New("Query myNode fail: " + err.Error())
	}

	if len(respBytes) > 0 {
		m.logger.Debug("Query myNode", "resp", string(respBytes))

		var response struct {
			Result int      `json:"_result"`
			Desc   string   `json:"_desc"`
			Data   NodeInfo `json:"data"`
		}

		err := json.Unmarshal(respBytes, &response)
		if err != nil {
			return "", "", "", -1, "", errors.New("Query myNode fail: " + err.Error())
		}

		if response.Result == 0 {
			// return myNode.NodeID, myNode.Owner.Encode(), myNode.Wallet.Encode(), myNode.Registered.Int64(), myNode.NodeType.Int64(), nil
			return response.Data.NodeID, response.Data.PublicKey, response.Data.Principal, int64(response.Data.Status), response.Data.NodeType, nil
		} else {
			return "", "", "", -1, "", errors.New("Query myNode fail: " + err.Error())
		}
	}
	return "", "", "", -1, "", errors.New("Query myNode fail")
}

func (m *MinerHubAgent) MyStack(nodeId string) (uint64, uint64, uint64, error) {
	//methodName := "myStake"
	//
	//var argType []idl.Type
	//argType = append(argType, new(idl.Text))
	//
	//argValue := []interface{}{
	//	nodeId,
	//}
	//arg, _ := idl.Encode(argType, argValue)
	//m.logger.Debug("MyStack", "argType", argType, "argValue", argValue)
	//
	//types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	//if err != nil {
	//	return 0, 0, 10000, err
	//}
	//m.logger.Debug("MyStack", "types", types, "result", result)
	////ouput-> types: [interface vec interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}]
	////ouput-> result: [[map[3054210041:[57 248 44 186 16 145 100 93 182 123 49 153 147 45 214 150 45 224 237 216 84 142 130 10 172 82 193 48 2] 4104166786:1685512050331435992 4135997916:0 656559709:16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC 947296307:[245 50 153 79 90 148 3 179 181 210 38 205 150 98 51 71 55 221 150 24 248 186 191 134 143 61 52 87 2]]]]
	//if result != nil && len(result) >= 3 {
	//	return result[0].(*big.Int).Uint64(), result[1].(*big.Int).Uint64(), result[2].(*big.Int).Uint64(), nil
	//}
	return 0, 0, 10000, nil
}

func (s *MinerHubAgent) getPrivateKey() *ecdsa.PrivateKey {
	networkPrivKey, err := s.secretsManager.GetSecret(secrets.ValidatorKey)
	if err != nil {
		return nil
	}

	decodedPrivKey, err := crypto.BytesToECDSAPrivateKey(networkPrivKey)
	if err != nil {
		return nil
	}

	return decodedPrivKey
}

// call EMCHub's registerNode api
func (m *MinerHubAgent) RegisterComputingNode(nodeId string, minerPrincipal string) error {
	privateKey := m.getPrivateKey()
	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		return errors.New("RegisterComputingNode fail: unable to extract key")
	}

	randnum := rand.Intn(1e6)
	message := address.String() + "," + minerPrincipal + "," + nodeId + "," + string(randnum)
	keccak256 := crypto.Keccak256([]byte(message))

	signature, err := crypto.Sign(
		privateKey,
		keccak256,
	)
	if err != nil {
		return errors.New("RegisterComputingNode fail: " + err.Error())
	}

	signatureHexString := hex.EncodeToString(signature)
	keccak256HexString := hex.EncodeToString(keccak256)

	m.logger.Info("RegisterComputingNode", "nodeId", nodeId, "public key", address.String(), "principal", minerPrincipal, "message", message, "keccak256", keccak256HexString, "signature", signatureHexString)

	var entity struct {
		NodeId    string `json:"nodeId"`
		NodeType  string `json:"nodeType"`
		PublicKey string `json:"publicKey"`
		Principal string `json:"principal"`
		Kecack256 string `json:"kecack256"`
		Signature string `json:"signature"`
	}
	entity.NodeId = nodeId
	entity.Principal = minerPrincipal
	entity.NodeType = "computing"
	entity.Signature = signatureHexString
	entity.Kecack256 = keccak256HexString
	entity.PublicKey = address.String()

	entityJsonBytes, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	respBytes, err := m.httpClient.SendPostJsonRequest(DEFAULT_HUB_HOST+"/api/v1/nodesign/add", entityJsonBytes)
	if err != nil {
		return errors.New("RegisterComputingNode fail: " + err.Error())
	}

	if len(respBytes) > 0 {
		m.logger.Debug("RegisterComputingNode", "resp", string(respBytes))

		var response struct {
			Result int    `json:"_result"`
			Desc   string `json:"_desc"`
		}
		err := json.Unmarshal(respBytes, &response)
		if err != nil {
			return errors.New("RegisterComputingNode fail: " + err.Error())
		}

		if response.Result == 0 {
			return nil
		} else {
			return errors.New("RegisterComputingNode fail: " + response.Desc)
		}
	}

	return errors.New("RegisterComputingNode fail")
}

func (m *MinerHubAgent) AddRouter(minerPrincipal string) error {

	return errors.New("AddRouter fail")
}

func (m *MinerHubAgent) RegisterValidatorNode(nodeId string, minerPrincipal string) error {

	return errors.New("RegisterNode fail")
}
func (m *MinerHubAgent) RegisterRouterNode(nodeId string, minerPrincipal string) error {

	return errors.New("RegisterNode fail")
}

// UnRegisterComputingNode
func (m *MinerHubAgent) UnRegisterComputingNode(nodeId string) error {
	privateKey := m.getPrivateKey()
	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		return errors.New("RegisterComputingNode fail: unable to extract key")
	}

	randnum := rand.Intn(1e6)
	message := address.String() + "," + nodeId + "," + string(randnum)
	keccak256 := crypto.Keccak256([]byte(message))

	signature, err := crypto.Sign(
		privateKey,
		keccak256,
	)
	if err != nil {
		return errors.New("RegisterComputingNode fail: " + err.Error())
	}

	signatureHexString := hex.EncodeToString(signature)
	keccak256HexString := hex.EncodeToString(keccak256)

	m.logger.Info("UnRegisterComputingNode", "nodeId", nodeId, "public key", address.String(), "message", message, "keccak256", keccak256HexString, "signature", signatureHexString)

	var entity struct {
		NodeId    string `json:"nodeId"`
		PublicKey string `json:"publicKey"`
		Kecack256 string `json:"kecack256"`
		Signature string `json:"signature"`
	}
	entity.NodeId = nodeId
	entity.Signature = signatureHexString
	entity.Kecack256 = keccak256HexString
	entity.PublicKey = address.String()

	entityJsonBytes, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	respBytes, err := m.httpClient.SendPostJsonRequest(DEFAULT_HUB_HOST+"/api/v1/nodesign/remove", entityJsonBytes)
	if err != nil {
		return errors.New("UnRegisterComputingNode fail: " + err.Error())
	}

	if len(respBytes) > 0 {
		m.logger.Debug("UnRegisterComputingNode", "resp", string(respBytes))

		var response struct {
			Result int    `json:"_result"`
			Desc   string `json:"_desc"`
		}
		err := json.Unmarshal(respBytes, &response)
		if err != nil {
			return errors.New("UnRegisterComputingNode fail: " + err.Error())
		}

		if response.Result == 0 {
			return nil
		} else {
			return errors.New("UnRegisterComputingNode fail: " + response.Desc)
		}
	}

	return errors.New("UnRegisterComputingNode fail")
}

// call miner canister's unRegisterNode method(update)
func (m *MinerHubAgent) UnregisterValidatorNode(nodeId string) error {

	return errors.New("UnregisterValidatorNode fail")
}

func (m *MinerHubAgent) UnRegisterRouterNode(nodeId string) error {

	return errors.New("UnregisterRouterNode fail")
}

func (m *MinerHubAgent) ListValidatorsNodeId() ([]string, error) {
	//nodeList := make([]string, 0)
	//methodName := "listValidatorNodes"
	//var argType []idl.Type
	//var argValue []interface{}
	//arg, _ := idl.Encode(argType, argValue)
	//
	//types, result, _, err := m.agent.Query(m.canister, methodName, arg)
	//if err != nil {
	//	m.logger.Error(err.Error())
	//	return nodeList, err
	//}
	////fmt.Println("types:", types[0].String(), "result:", result)
	////types: interface vec interface record {0:text; 1:interface record {656559709:text; 947296307:principal; 3054210041:principal; 4104166786:int; 4135997916:nat}} result: [[map[0:16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y 947296307:[36 168 117 208 89 72 247 191 251 207 187 65 155 215 201 21 37 102 213 51 157 96 202 221 113 185 26 84 2]]] map[0:16Uiu2HAmEoDReK7pKygYYYFgJ8uuXS8oWsYFWiiEbCSF9HjYcih2 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmEoDReK7pKygYYYFgJ8uuXS8oWsYFWiiEbCSF9HjYcih2 947296307:[4 147 74 93 122 138 141 1 188 183 152 85 46 153 244 27 237 124 74 239 66 112 148 29 10 21 174 46 2]]] map[0:16Uiu2HAkyPw8SEeDpErEwcEZ2QtXzPq5KQf4woWybsr7KN6VH7yX 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685612201271035791 4135997916:1 656559709:16Uiu2HAkyPw8SEeDpErEwcEZ2QtXzPq5KQf4woWybsr7KN6VH7yX 947296307:[137 93 57 21 145 14 171 17 20 8 27 60 176 216 139 147 110 198 119 130 170 1 175 38 61 19 108 211 2]]] map[0:16Uiu2HAmPfFVHNnYKdDQywJXnzbgM1MdAi6P1MsCkxN7Hr6VaiYa 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmPfFVHNnYKdDQywJXnzbgM1MdAi6P1MsCkxN7Hr6VaiYa 947296307:[42 174 237 43 160 163 242 37 48 118 200 166 112 189 56 255 19 59 90 30 53 188 133 153 30 253 236 40 2]]] map[0:16Uiu2HAmTPfBgUkQ4V8qaBvTaJp54Cm32TWGvYZaxcuPxoaSbZAS 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmTPfBgUkQ4V8qaBvTaJp54Cm32TWGvYZaxcuPxoaSbZAS 947296307:[102 242 43 171 69 192 233 201 175 42 233 72 161 61 119 192 214 32 225 200 111 96 26 111 57 76 242 248 2]]] map[0:16Uiu2HAmGpKZdnpaaYgKTZqagLVJcnMphdeqaHtKBaFFkb5MYRUy 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605846436938574 4135997916:1 656559709:16Uiu2HAmGpKZdnpaaYgKTZqagLVJcnMphdeqaHtKBaFFkb5MYRUy 947296307:[68 196 3 132 35 19 140 71 180 180 247 198 72 205 39 164 52 180 106 88 60 107 211 100 131 101 160 103 2]]] map[0:16Uiu2HAm7BqtmjH7JECa5Y4iNgiZXuet3HqXYZbeXNN9XiQwgSbf 1:map[3054210041:[128 69 101 81 144 137 246 95 120 89 13 116 173 114 174 17 9 171 108 77 186 206 134 20 204 76 189 105 2] 4104166786:1685605847289287331 4135997916:1 656559709:16Uiu2HAm7BqtmjH7JECa5Y4iNgiZXuet3HqXYZbeXNN9XiQwgSbf 947296307:[46 113 59 128 101 160 251 203 212 119 205 85 133 9 206 253 168 134 11 30 174 45 49 197 147 0 246 137 2]]]]]
	//m.logger.Debug("ListValidators", "types", types[0].String(), "result", result)
	////fmt.Println("len(result):", len(result))
	//if result != nil && len(result) > 0 {
	//	vec := result[0].([]interface{})
	//	for _, recordVec := range vec {
	//		record := recordVec.(map[string]interface{})
	//
	//		validatorNode := ValidorNodeInfo{}
	//		utils.Decode(&validatorNode, record)
	//
	//		nodeList = append(nodeList, validatorNode.NodeID)
	//	}
	//	return nodeList, nil
	//}
	return nil, nil
}
