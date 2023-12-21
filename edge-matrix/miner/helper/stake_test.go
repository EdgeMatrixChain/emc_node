package helper

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func Test_GetLockedAmount(t *testing.T) {
	client, err := ethclient.Dial("https://arbitrum-goerli.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0xDC1E36492317D1A79c6e7DfA772e0D91930d99ea")
	instance, err := NewStake(address, client)
	if err != nil {
		log.Fatal(err)
	}

	amount, err := instance.GetLockedAmount(&bind.CallOpts{
		Pending: false,
	}, common.HexToAddress("0xd5e1c4e65860e7131082b14799d6251a9a33a163"))
	if err != nil {
		log.Fatal(err)
	}

	t.Log("amount:", amount)
}
