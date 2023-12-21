package helper

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"testing"
)

func Test_ChainID(t *testing.T) {
	id, err := ChainID("https://arbitrum-goerli.publicnode.com")
	if err != nil {
		panic(err)
	}
	t.Log("id:", id)
}

func Test_Call(t *testing.T) {
	bytes, err := Call(
		"https://arbitrum-goerli.publicnode.com",
		"0xDC1E36492317D1A79c6e7DfA772e0D91930d99ea",
		"0x70a08231000000000000000000000000d5e1c4e65860e7131082b14799d6251a9a33a163")
	if err != nil {
		panic(err)
	}

	t.Log("resp:", hexutil.Encode(bytes))
}
