package helper

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func Test_GetNodeStakeAmount(t *testing.T) {
	client, err := ethclient.Dial("https://sepolia-rollup.arbitrum.io/rpc")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0xbfbc3BF85FBA818fc49A0354D2C84623cE711b63")
	instance, err := NewStake(address, client)
	if err != nil {
		log.Fatal(err)
	}

	amount, err := instance.BalanceOfNode(&bind.CallOpts{
		Pending: false,
	}, "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")
	if err != nil {
		log.Fatal(err)
	}

	t.Log("amount:", amount)
}
