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
	//privateKey := m.getPrivateKey()
	//address, err := crypto.GetAddressFromKey(privateKey)
	//if err != nil {
	//	return errors.New("RegisterComputingNode fail: unable to extract key")
	//}
	//
	//randnum := rand.Intn(1e6)
	//message := address.String() + "," + nodeId + "," + string(randnum)
	//keccak256 := crypto.Keccak256([]byte(message))
	//
	//signature, err := crypto.Sign(
	//	privateKey,
	//	keccak256,
	//)
	//if err != nil {
	//	return errors.New("RegisterComputingNode fail: " + err.Error())
	//}
	//
	//signatureHexString := hex.EncodeToString(signature)
	//keccak256HexString := hex.EncodeToString(keccak256)
	//
	//m.logger.Info("UnRegisterComputingNode", "nodeId", nodeId, "public key", address.String(), "message", message, "keccak256", keccak256HexString, "signature", signatureHexString)
	//
	//var entity struct {
	//	NodeId    string `json:"nodeId"`
	//	PublicKey string `json:"publicKey"`
	//	Kecack256 string `json:"kecack256"`
	//	Signature string `json:"signature"`
	//}
	//entity.NodeId = nodeId
	//entity.Signature = signatureHexString
	//entity.Kecack256 = keccak256HexString
	//entity.PublicKey = address.String()
	//
	//entityJsonBytes, err := json.Marshal(entity)
	//if err != nil {
	//	return err
	//}
	//respBytes, err := m.httpClient.SendPostJsonRequest(DEFAULT_HUB_HOST+"/api/v1/nodesign/remove", entityJsonBytes)
	//if err != nil {
	//	return errors.New("UnRegisterComputingNode fail: " + err.Error())
	//}
	//
	//if len(respBytes) > 0 {
	//	m.logger.Debug("UnRegisterComputingNode", "resp", string(respBytes))
	//
	//	var response struct {
	//		Result int    `json:"_result"`
	//		Desc   string `json:"_desc"`
	//	}
	//	err := json.Unmarshal(respBytes, &response)
	//	if err != nil {
	//		return errors.New("UnRegisterComputingNode fail: " + err.Error())
	//	}
	//
	//	if response.Result == 0 {
	//		return nil
	//	} else {
	//		return errors.New("UnRegisterComputingNode fail: " + response.Desc)
	//	}
	//}

	return errors.New("UnRegisterComputingNode fail")
}

// call miner canister's unRegisterNode method(update)
func (m *MinerHubAgent) UnregisterValidatorNode(nodeId string) error {

	return errors.New("UnregisterValidatorNode fail")
}

func (m *MinerHubAgent) UnRegisterRouterNode(nodeId string) error {

	return errors.New("UnregisterRouterNode fail")
}
