package rpc

import (
	"encoding/base64"
	"fmt"
	"github.com/emc-protocol/edge-matrix/crypto"
	"log"
	"testing"
)

func TestGetNextNonce(t *testing.T) {
	client := NewDefaultJsonRpcClient()
	address := "0x0aF137aa3EcC7d10d926013ee34049AfA77382e6"
	nonce, err := client.GetNextNonce(address)
	if err != nil {
		t.Error(err)
	}
	t.Log("nonce:", nonce)
}

func TestSendPocCpu(t *testing.T) {
	seed := "0x4fe0d01ce75cb08562e1c67a6592633eaa2ffb82b8b334efa76f3adad77bea67"
	nodeId := "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC"
	inputString := fmt.Sprintf("{\"peerId\": \"%s\",\"endpoint\": \"/poc_cpu\",\"Input\": {\"seed\": \"%s\"}}", nodeId, seed)

	client := NewDefaultJsonRpcClient()
	privateKey, err := crypto.BytesToECDSAPrivateKey([]byte("d4ffa0ca147fce3cacbffebf0c411010bd0c8e5a27f16c918032f4ddd5c2665a"))
	if err != nil {
		t.Error(err)
	}
	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		t.Fatalf("unable to extract key, %v", err)
	}

	nonce, err := client.GetNextNonce(address.String())
	if err != nil {
		t.Error(err)
	}

	response, err := client.SendRawTelegram(
		EdgeCallPrecompile,
		nonce,
		inputString,
		privateKey,
	)
	if err != nil {
		t.Error(err)
	}

	t.Log("TelegramHash:", response.Result.TelegramHash)
	decodeBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		log.Fatalln(err)
	}
	t.Log("Response:", string(decodeBytes))
}

func TestSendPocRequest(t *testing.T) {
	validatorNodeID := "16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y"
	nodeId := "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC"

	input := fmt.Sprintf("{\"peerId\": \"%s\",\"endpoint\": \"/poc_request\",\"Input\": {\"node_id\": \"%s\"}}", validatorNodeID, nodeId)
	client := NewDefaultJsonRpcClient()
	privateKey, err := crypto.BytesToECDSAPrivateKey([]byte("d4ffa0ca147fce3cacbffebf0c411010bd0c8e5a27f16c918032f4ddd5c2665a"))
	if err != nil {
		t.Error(err)
	}
	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		t.Fatalf("unable to extract key, %v", err)
	}

	nonce, err := client.GetNextNonce(address.String())
	if err != nil {
		t.Error(err)
	}

	response, err := client.SendRawTelegram(
		EdgeCallPrecompile,
		nonce,
		input,
		privateKey,
	)
	if err != nil {
		t.Error(err)
	}

	t.Log("TelegramHash:", response.Result.TelegramHash)
	decodeBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		log.Fatalln(err)
	}
	t.Log("Response:", string(decodeBytes))
}
func TestCallInfo(t *testing.T) {
	client := NewDefaultJsonRpcClient()
	privateKey, err := crypto.BytesToECDSAPrivateKey([]byte("d4ffa0ca147fce3cacbffebf0c411010bd0c8e5a27f16c918032f4ddd5c2665a"))
	if err != nil {
		t.Error(err)
	}
	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		t.Fatalf("unable to extract key, %v", err)
	}

	nonce, err := client.GetNextNonce(address.String())
	if err != nil {
		t.Error(err)
	}
	//nonce = 33
	t.Log("nonce:", nonce)
	input := `{"peerId": "16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y","endpoint": "/info","input": {}}`
	response, err := client.SendRawTelegram(
		EdgeCallPrecompile,
		nonce,
		input,
		privateKey,
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log("TelegramHash:", response.Result.TelegramHash)
	decodeBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		log.Fatalln(err.Error())
	}
	t.Log("Response:", string(decodeBytes))
}

func TestCallPocCpuRequest(t *testing.T) {
	client := NewDefaultJsonRpcClient()
	privateKey, err := crypto.BytesToECDSAPrivateKey([]byte("d4ffa0ca147fce3cacbffebf0c411010bd0c8e5a27f16c918032f4ddd5c2665a"))
	if err != nil {
		t.Error(err)
	}
	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		t.Fatalf("unable to extract key, %v", err)
	}

	nonce, err := client.GetNextNonce(address.String())
	if err != nil {
		t.Error(err)
	}
	//nonce = 33
	t.Log("nonce:", nonce)
	input := `{"peerId": "16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y","endpoint": "/poc_cpu_request","input": {"node_id" : "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC"}}`
	response, err := client.SendRawTelegram(
		EdgeCallPrecompile,
		nonce,
		input,
		privateKey,
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log("TelegramHash:", response.Result.TelegramHash)
	decodeBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
	if err != nil {
		log.Fatalln(err.Error())
	}
	t.Log("Response:", string(decodeBytes))
}
