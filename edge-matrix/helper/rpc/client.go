package rpc

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/types"
	"math/big"
)

const DEFAULT_RPC_URL = "https://oregon.edgematrix.xyz"
const TESTNET_ID = 2
const MAINET_ID = 1

var (
	// EdgeSubscribeRegisterPrecompile is and address of edge subscribe register precompile
	EdgeSubscribeRegisterPrecompile = types.StringToAddress("0x3000")
	// EdgeCallPrecompile is and address of edge call precompile
	EdgeCallPrecompile = types.StringToAddress("0x3001")
	// EdgeRtcSubjectPrecompile is and address of edge subject precompile
	EdgeRtcSubjectPrecompile = types.StringToAddress("0x3101")
)

type JsonRpcClient struct {
	httpClient *FastHttpClient
	signer     *crypto.EIP155Signer
	rpcUrl     string
}

type StringResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Result  string `json:"result"`
	Error   Error  `json:"error"`
}

type TelegramResponse struct {
	JsonRpc string `json:"jsonrpc"`
	Id      uint64 `json:"id"`
	Result  Reuslt `json:"result"`
	Error   Error  `json:"error"`
}

type Reuslt struct {
	// "{\"telegram_hash\":\"__HashHexString__\",\"response\":\"__Base64String__\"}"
	TelegramHash string `json:"telegram_hash"`
	Response     string `json:"response"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RawResponse struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      uint64          `json:"id"`
	Result  json.RawMessage `json:"result"`
}

func NewDefaultJsonRpcClient() *JsonRpcClient {
	return &JsonRpcClient{
		httpClient: NewDefaultHttpClient(),
		signer:     crypto.NewEIP155Signer(chain.AllForksEnabled.At(0), big.NewInt(TESTNET_ID).Uint64()),
		rpcUrl:     DEFAULT_RPC_URL,
	}
}

// Returns next nonce value for address
func (c *JsonRpcClient) GetNextNonce(address string) (uint64, error) {
	postJson := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"edge_getTelegramCount\",\"params\":[\"%s\"],\"id\":1}", address)
	bytes, err := c.httpClient.SendPostJsonRequest(c.rpcUrl, []byte(postJson))
	if err != nil {
		return 0, err
	}
	response := &StringResponse{}
	err = json.Unmarshal(bytes, response)
	if err != nil {
		return 0, err
	}
	decodeUint64, err := hex.DecodeUint64(response.Result)
	if err != nil {
		return 0, err
	}
	return decodeUint64, nil
}

// Call edge_sendRawTelegram api method
func (c *JsonRpcClient) SendRawTelegram(to types.Address, nonce uint64, input string, privateKey *ecdsa.PrivateKey) (*TelegramResponse, error) {
	//input := `{"peerId": "16Uiu2HAkw3hzExAr4CXDrC2VQfeXgrJdn5bVaV58XRbXad77EU9V","endpoint": "/echo","input": "hello"}`
	tele := &types.Telegram{
		To:       &to,
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(0),
		Nonce:    nonce,
		Input:    []byte(input),
	}
	_, err := json.Marshal(tele)
	if err != nil {
		return nil, errors.New("json.Marshal err: " + err.Error())
	}

	signedTx, signErr := c.signer.SignTele(tele, privateKey)
	if signErr != nil {
		return nil, errors.New("Unable to sign transaction")
	}

	bytes := signedTx.MarshalRLP()
	postJson := fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"method\":\"edge_sendRawTelegram\",\"params\":[\"%s\"],\"id\":1}", hex.EncodeToHex(bytes))
	bytes, err = c.httpClient.SendPostJsonRequest(c.rpcUrl, []byte(postJson))
	if err != nil {
		return nil, errors.New("SendPostJsonRequest err:" + err.Error())
	}
	// {"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"nonce too low"}}
	//fmt.Println("SendRawTelegram bytes:", string(bytes))

	sresp := &StringResponse{}
	err = json.Unmarshal(bytes, sresp)
	if err != nil {
		return nil, errors.New("json.Unmarshal sresp err:" + err.Error())
	}
	if sresp.Error.Code < 0 {
		return nil, errors.New("sresp err:" + sresp.Error.Message)
	}
	response := &TelegramResponse{
		JsonRpc: sresp.JsonRpc,
		Id:      sresp.Id,
	}
	r := &Reuslt{}
	err = json.Unmarshal([]byte(sresp.Result), r)
	if err != nil {
		response.Error = Error{Code: -1, Message: "json.Unmarshal Result err:" + err.Error() + " sresp.Result: " + sresp.Result}
	}
	response.Result = *r
	return response, nil
}
