package helper

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func ChainID(rpcUrl string) (*big.Int, error) {
	ec, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	defer ec.Close()

	id, err := ec.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	return id, nil
}

func Call(rpcUrl string, contractAddress string, dataHexString string) ([]byte, error) {
	ec, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}
	defer ec.Close()

	contractAddr := common.HexToAddress(contractAddress)
	callMsg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: common.FromHex(dataHexString),
	}

	resp, err := ec.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Error calling contract: %v", err)
	}

	return resp, nil
}
